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
use std::collections::VecDeque;
use std::sync::mpsc;
use std::sync::Arc;
use std::sync::Mutex;

// Global state for run/stop button
static IS_RUNNING: AtomicBool = AtomicBool::new(false);

// Global variable to store tail process
static mut TAIL_PROCESS: Option<Child> = None;

// Add a constant for tracking script modification
static IS_SCRIPT_MODIFIED: AtomicBool = AtomicBool::new(false);

// Default validator script content
const DEFAULT_SCRIPT: &str = r#"#!/bin/bash
exec PATH_OF_solana-validator \
    --identity PATH_OF_identity.json \
    --vote-account PUBLIC_ADDRESS_OF_VOTE \
    --known-validator Abt4r6uhFs7yPwR3jT5qbnLjBtasgHkRVAd1W6H5yonT \
    --known-validator 5NfpgFCwrYzcgJkda9bRJvccycLUo3dvVQsVAK2W43Um \
    --only-known-rpc \
    --log validator.log \
    --ledger ledger \
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
    &"#;

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

// Add this function to check if validator is running
fn is_validator_running() -> bool {
    let output = Command::new("ps")
        .args(["aux"])
        .output()
        .map(|output| {
            let processes = String::from_utf8_lossy(&output.stdout);
            processes.contains("solana-validator")
        })
        .unwrap_or(false);
    
    output
}

// Extract validator path from script content
fn extract_validator_path(script_content: &str) -> Option<String> {
    let re = Regex::new(r"exec\s+([^\s\\]+)").ok()?;
    re.captures(script_content)
        .and_then(|caps| caps.get(1))
        .map(|m| m.as_str().to_string())
}

// Add function to extract ledger path from script content
fn extract_ledger_path(script_content: &str) -> Option<String> {
    let re = Regex::new(r"--ledger\s+([^\s\\]+)").ok()?;
    re.captures(script_content)
        .and_then(|caps| caps.get(1))
        .map(|m| m.as_str().to_string())
}

// Create and return the validator view layout
pub fn get_validator_view() -> LinearLayout {
    // Check initial validator state
    let is_running = is_validator_running();
    IS_RUNNING.store(is_running, Ordering::SeqCst);

    // Create dashboard with status information
    let dashboard = Panel::new(
        LinearLayout::horizontal()
            .child(TextView::new("Validator Status: "))
            .child(TextView::new(if is_running { "RUNNING" } else { "STOPPED" }))
    )
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
        .child(Button::new("Save Script", |s| {
            save_script(s);
        }))
        .child(DummyView.fixed_width(4))
        .child(Button::new(if is_running { "Stop Validator" } else { "Start Validator" }, move |s| {
            // Auto save if modified before running
            if IS_SCRIPT_MODIFIED.load(Ordering::SeqCst) {
                save_script(s);
            }
            toggle_run_stop(s);
        }).with_name("run_button"))
        .child(DummyView.fixed_width(4))
        .child(Button::new("Refresh Logs", move |s| {
            if is_validator_running() {
                if let Some(log_path) = extract_log_path(&get_script_content()) {
                    // Print log path to logs area
                    update_logs(s, &format!("Monitoring log file: {}", &log_path));
                    
                    // Start log monitoring
                    if let Some(process) = start_log_monitor(s, &log_path) {
                        unsafe {
                            TAIL_PROCESS = Some(process);
                        }
                    }
                }
            }
        }));

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
        // Trim any trailing whitespace or newlines
        view.get_content().trim_end().to_string()
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
                    
                    IS_SCRIPT_MODIFIED.store(false, Ordering::SeqCst);
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
    let log_path = re.captures(script_content)
        .and_then(|caps| caps.get(1))
        .map(|m| m.as_str().to_string())?;

    // Convert relative path to absolute path if needed
    if let Ok(exe_path) = env::current_exe() {
        if let Some(exe_dir) = exe_path.parent() {
            let absolute_path = exe_dir.join(&log_path);
            return Some(absolute_path.to_string_lossy().into_owned());
        }
    }
    Some(log_path)
}

// Add a constant for maximum log lines
const MAX_LOG_LINES: usize = 100;

// Start tail process and monitor its output
fn start_log_monitor(siv: &mut Cursive, log_path: &str) -> Option<Child> {
    // Create log file if it doesn't exist
    if !std::path::Path::new(log_path).exists() {
        if let Err(e) = std::fs::File::create(log_path) {
            update_logs(siv, &format!("Failed to create log file: {}", e));
            return None;
        }
    }
    
    let mut cmd = Command::new("tail")
        .args(["-f", "-n", "10", log_path])
        .stdout(std::process::Stdio::piped())
        .spawn()
        .ok()?;
    
    let stdout = cmd.stdout.take()?;
    let reader = BufReader::new(stdout);
    let siv = siv.cb_sink().clone();
    
    // Create a channel for log messages
    let (tx, rx) = mpsc::channel();
    let tx_clone = tx.clone();

    // Spawn a thread to read logs
    std::thread::spawn(move || {
        for line in reader.lines() {
            if let Ok(line) = line {
                if tx_clone.send(line).is_err() {
                    break;
                }
            }
        }
    });

    // Spawn another thread to batch process logs
    std::thread::spawn(move || {
        let log_buffer = Arc::new(Mutex::new(VecDeque::with_capacity(MAX_LOG_LINES)));
        let mut batch = Vec::new();
        let mut last_update = std::time::Instant::now();

        while let Ok(line) = rx.recv() {
            batch.push(line);

            // Update UI if we have collected enough lines or enough time has passed
            if batch.len() >= 10 || last_update.elapsed() >= std::time::Duration::from_millis(100) {
                if !batch.is_empty() {
                    let messages = batch.join("\n");
                    let buffer_clone = Arc::clone(&log_buffer);
                    let _ = siv.send(Box::new(move |s| {
                        update_logs_batch(s, &messages, &buffer_clone);
                    }));
                    batch.clear();
                    last_update = std::time::Instant::now();
                }
            }
        }
    });

    Some(cmd)
}

// Update logs with batched messages
fn update_logs_batch(siv: &mut Cursive, messages: &str, log_buffer: &Arc<Mutex<VecDeque<String>>>) {
    // Clean ANSI escape sequences
    let clean_messages = clean_log_message(messages);
    
    // Get lock on buffer
    if let Ok(mut buffer) = log_buffer.lock() {
        // Split messages into lines and add to buffer
        for line in clean_messages.lines() {
            buffer.push_back(line.to_string());
            // Keep buffer size limited
            while buffer.len() > MAX_LOG_LINES {
                buffer.pop_front();
            }
        }

        // Update TextView with all buffered logs
        siv.call_on_name("log_view", |view: &mut Panel<ScrollView<TextView>>| {
            let text_view = view.get_inner_mut().get_inner_mut();
            let content = buffer.iter().cloned().collect::<Vec<_>>().join("\n");
            text_view.set_content(content);
        });
    }
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
            // Start the validator script
            if let Ok(exe_path) = env::current_exe() {
                if let Some(exe_dir) = exe_path.parent() {
                    let script_path = exe_dir.join("validator-testnet.sh");
                    
                    // Execute the script
                    match Command::new("bash")
                        .arg(&script_path)
                        .spawn() {
                        Ok(_) => {
                            update_logs(siv, "Validator script started successfully!");
                        },
                        Err(e) => {
                            update_logs(siv, &format!("Failed to start validator: {}", e));
                            return;
                        }
                    }
                }
            }

            // Start log monitoring
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
        // Get script content
        let script_content = siv.call_on_name("script_content", |view: &mut TextArea| {
            view.get_content().to_string()
        }).unwrap_or_default();

        // Get validator path and ledger path
        if let (Some(validator_path), Some(ledger_path)) = (
            extract_validator_path(&script_content),
            extract_ledger_path(&script_content)
        ) {
            // Execute solana-validator exit command with ledger path
            match Command::new(&validator_path)
                .args(["--ledger", &ledger_path, "exit", "-f"])
                .status() {
                Ok(_) => {
                    update_logs(siv, "Validator stopping gracefully...");
                },
                Err(e) => {
                    update_logs(siv, &format!("Failed to stop validator: {}", e));
                    return;
                }
            }
        } else {
            update_logs(siv, "Could not find validator path or ledger path in script");
            return;
        }

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