use std::process::{Command, Child};
use std::io::{BufRead, BufReader};
use std::sync::atomic::{AtomicBool, Ordering};
use std::env;
use std::fs;
use regex::Regex;
use cursive::views::{LinearLayout, Panel, TextView, TextArea, Button, DummyView, ResizedView, ScrollView};
use cursive::traits::*;
use cursive::Cursive;
use lazy_static::lazy_static;

// Global state for run/stop button
static IS_RUNNING: AtomicBool = AtomicBool::new(false);

// Global variable to store tail process
static mut TAIL_PROCESS: Option<Child> = None;

// Default validator script content
const DEFAULT_SCRIPT: &str = r#"#!/bin/bash
exec PATH_OF_SOLANA_VALIDATOR \
    --identity PATH_OF_identity.json \
    --vote-account PUBLIC_ADDRESS_OF_vote.json \
    --known-validator Abt4r6uhFs7yPwR3jT5qbnLjBtasgHkRVAd1W6H5yonT \
    --known-validator 5NfpgFCwrYzcgJkda9bRJvccycLUo3dvVQsVAK2W43Um \
    --only-known-rpc \
    --log PATH_OF_validator.log \
    --ledger FOLDER_PATH_OF_ledger \
    --rpc-port 8899 \
    --full-rpc-api \
    --dynamic-port-range 8000-8020 \
    --entrypoint xolana.xen.network:8001 \
    --entrypoint owlnet.dev:8001 \
    --wal-recovery-mode skip_any_corrupted_record \
    --limit-ledger-size 50000000 \
    --enable-rpc-transaction-history \
    --enable-extended-tx-metadata-storage \
    --rpc-pubsub-enable-block-subscription \
    --full-snapshot-interval-slots 5000 \
    --maximum-incremental-snapshots-to-retain 10 \
    --maximum-full-snapshots-to-retain 50 \
    &
"#;

// Initialize regex pattern for ANSI escape codes
lazy_static! {
    static ref ANSI_ESCAPE_RE: Regex = Regex::new(r"\x1B\[[0-9;]*[a-zA-Z]|\x1B\[[0-9;]*m").unwrap();
}

// Function to get script content
fn get_script_content() -> String {
    // Get current executable path
    if let Ok(exe_path) = env::current_exe() {
        // Get the directory containing the executable
        if let Some(exe_dir) = exe_path.parent() {
            // Create script path in the same directory
            let script_path = exe_dir.join("validator-testnet.sh");
            
            // Try to read existing script
            if let Ok(content) = fs::read_to_string(&script_path) {
                return content;
            }
        }
    }
    
    // Return default script if file doesn't exist or can't be read
    DEFAULT_SCRIPT.to_string()
}

// Create and return the validator view layout
pub fn get_validator_view() -> LinearLayout {
    // Create three sections
    let dashboard = Panel::new(TextView::new("Validator Dashboard"))
        .title("Dashboard")
        .full_width()
        .fixed_height(5);

    // Create config section with TextArea and buttons
    let text_area = TextArea::new()
        .content(get_script_content())
        .with_name("script_content")
        .min_height(10)
        .max_height(24);

    // Create button layout with space between buttons
    let button_layout = LinearLayout::horizontal()
        .child(Button::new("Save", |s| {
            save_script(s);
        }))
        .child(DummyView.fixed_width(4))
        .child(Button::new("Run", move |s| {
            toggle_run_stop(s);
        }).with_name("run_button"));

    let config_content = LinearLayout::vertical()
        .child(ResizedView::with_full_screen(text_area))
        .child(button_layout);

    let config = Panel::new(config_content)
        .title("Config")
        .full_width()
        .min_height(10);

    let logs = Panel::new(
        ScrollView::new(TextView::new(""))
            .scroll_strategy(cursive::view::ScrollStrategy::StickToBottom)
    )
    .title("Logs")
    .with_name("log_view")
    .full_width()
    .min_height(8);

    // Combine sections vertically
    let layout = LinearLayout::vertical()
        .child(dashboard)
        .child(config)
        .child(logs);
    
    layout
}

// Function to save script content to file
fn save_script(siv: &mut cursive::Cursive) {
    let content = siv.call_on_name("script_content", |view: &mut TextArea| {
        view.get_content().to_string()
    }).unwrap_or_default();
    
    // Get current executable path
    if let Ok(exe_path) = env::current_exe() {
        // Get the directory containing the executable
        if let Some(exe_dir) = exe_path.parent() {
            // Create script path in the same directory
            let script_path = exe_dir.join("validator-testnet.sh");
            
            // Save the content to file
            match fs::write(&script_path, content) {
                Ok(_) => {
                    // Make the script executable (Unix-like systems only)
                    #[cfg(unix)]
                    {
                        use std::os::unix::fs::PermissionsExt;
                        if let Ok(metadata) = fs::metadata(&script_path) {
                            let mut perms = metadata.permissions();
                            perms.set_mode(0o755); // rwxr-xr-x
                            let _ = fs::set_permissions(&script_path, perms);
                        }
                    }
                    
                    // Update log to show success
                    update_logs(siv, "Script validator-testnet.sh saved successfully!");
                },
                Err(e) => {
                    // Update log to show error
                    update_logs(siv, &format!("Failed to save script: {}", e));
                }
            }
        }
    }
}

// Extract log file path from script content using regex
fn extract_log_path(script_content: &str) -> Option<String> {
    let re = Regex::new(r"--log\s+([^\s\\]+)").ok()?;
    re.captures(script_content)
        .and_then(|caps| caps.get(1))
        .map(|m| m.as_str().to_string())
}

// Start tail process and monitor its output
fn start_log_monitor(siv: &mut Cursive, log_path: &str) -> Option<Child> {
    // Start tail command with -f (follow) and -n 10 (last 10 lines)
    // Added --retry to keep trying if the file is inaccessible
    let mut cmd = Command::new("tail")
        .args(["-f", "-n", "10", "--retry", log_path])
        .stdout(std::process::Stdio::piped())
        .spawn()
        .ok()?;
    
    // Get stdout handle from the process
    let stdout = cmd.stdout.take()?;
    let reader = BufReader::new(stdout);
    let siv = siv.cb_sink().clone();

    // Spawn a new thread to read output
    std::thread::spawn(move || {
        for line in reader.lines() {
            if let Ok(line) = line {
                // Send each line to the UI
                let _ = siv.send(Box::new(move |s| {
                    update_logs(s, &line);
                }));
            }
        }
    });

    Some(cmd)
}

// Toggle between Run and Stop states
fn toggle_run_stop(siv: &mut Cursive) {
    let is_running = IS_RUNNING.load(Ordering::SeqCst);
    
    if !is_running {
        // Get script content from TextArea
        let script_content = siv.call_on_name("script_content", |view: &mut TextArea| {
            view.get_content().to_string()
        }).unwrap_or_default();
        
        // Extract log path and start monitoring
        if let Some(log_path) = extract_log_path(&script_content) {
            if let Some(process) = start_log_monitor(siv, &log_path) {
                unsafe {
                    TAIL_PROCESS = Some(process);
                }
            }
        }
        
        // Update button state to "Stop"
        siv.call_on_name("run_button", |button: &mut Button| {
            button.set_label("Stop");
        });
        IS_RUNNING.store(true, Ordering::SeqCst);
        
    } else {
        // Stop log monitoring by killing tail process
        unsafe {
            if let Some(mut process) = TAIL_PROCESS.take() {
                let _ = process.kill();
                let _ = process.wait();
            }
        }
        
        // Update button state back to "Run"
        siv.call_on_name("run_button", |button: &mut Button| {
            button.set_label("Run");
        });
        IS_RUNNING.store(false, Ordering::SeqCst);
    }
}

// Clean ANSI escape sequences from log message
fn clean_log_message(message: &str) -> String {
    ANSI_ESCAPE_RE.replace_all(message, "").to_string()
}

// Update the logs panel with new content
fn update_logs(siv: &mut Cursive, message: &str) {
    // Clean ANSI escape sequences before displaying
    let clean_message = clean_log_message(message);
    
    siv.call_on_name("log_view", |view: &mut Panel<ScrollView<TextView>>| {
        view.get_inner_mut().get_inner_mut().append(&clean_message);
        view.get_inner_mut().get_inner_mut().append("\n");
    });
}