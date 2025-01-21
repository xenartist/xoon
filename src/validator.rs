use std::process::{Command, Child, Stdio};
use std::io::{BufRead, BufReader};
use std::sync::atomic::{AtomicBool, Ordering};
use std::env;
use std::fs::{self};
use regex::Regex;
use cursive::views::{LinearLayout, Panel, TextView, TextArea, Button, DummyView, ResizedView, ScrollView, Dialog, RadioGroup};
use cursive::traits::*;
use cursive::Cursive;
use lazy_static::lazy_static;
use std::collections::VecDeque;
use std::sync::mpsc;
use std::sync::Arc;
use std::sync::Mutex;
use cursive::theme::{BaseColor, Color, Style};
use cursive::utils::markup::StyledString;
use std::path::PathBuf;


// Add a constant for tracking script modification
static IS_SCRIPT_MODIFIED: AtomicBool = AtomicBool::new(false);

// Add a constant for tracking auto check status
static IS_AUTO_CHECKING: AtomicBool = AtomicBool::new(false);

// Default validator script content
const DEFAULT_TESTNET_SCRIPT: &str = r#"#!/bin/bash
exec solana-validator \
    --identity ~/.config/solana/identity.json \
    --vote-account ~/.config/solana/vote.json \
    --known-validator Abt4r6uhFs7yPwR3jT5qbnLjBtasgHkRVAd1W6H5yonT \
    --known-validator 5NfpgFCwrYzcgJkda9bRJvccycLUo3dvVQsVAK2W43Um \
    --known-validator FcrZRBfVk2h634L9yvkysJdmvdAprq1NM4u263NuR6LC \
    --known-validator Tpsu5EYTJAXAat19VEh54zuauHvUBuryivSFRC3RiFk \
    --only-known-rpc \
    --log ./validator.log \
    --ledger ./ledger \
    --rpc-port 8899 \
    --full-rpc-api \
    --dynamic-port-range 8000-8020 \
    --entrypoint testnet.x1.xyz:8001 \
    --entrypoint testnet.x1.xyz:8000 \
    --entrypoint testnet.x1.xyz:8000 \
    --entrypoint owlnet.dev:8001 \
    --wal-recovery-mode skip_any_corrupted_record \
    --limit-ledger-size 50000000 \
    --enable-rpc-transaction-history \
    --enable-extended-tx-metadata-storage \
    --rpc-pubsub-enable-block-subscription \
    --full-snapshot-interval-slots 5000 \
    --maximum-incremental-snapshots-to-retain 10 \
    --maximum-full-snapshots-to-retain 50"#;

const DEFAULT_MAINNET_SCRIPT: &str = r#"#!/bin/bash
# Mainnet validator script will be added here
"#;

// Initialize regex pattern for ANSI escape codes
lazy_static! {
    static ref ANSI_ESCAPE_RE: Regex = Regex::new(r"\x1B\[[0-9;]*[a-zA-Z]|\x1B\[[0-9;]*m").unwrap();
}

lazy_static! {
    static ref CURRENT_NETWORK: Mutex<String> = Mutex::new("testnet".to_string());
}

// Function to get script content
fn get_script_content(network: &str) -> String {
    // Get current executable path
    if let Ok(exe_path) = env::current_exe() {
        // Get the directory containing the executable
        if let Some(exe_dir) = exe_path.parent() {
            // Create script path based on network
            let script_name = if network == "mainnet" {
                "validator-mainnet.sh"
            } else {
                "validator-testnet.sh"
            };
            let script_path = exe_dir.join(script_name);
            
            // Try to read existing script
            if let Ok(content) = fs::read_to_string(&script_path) {
                return content;
            }
        }
    }
    
    // Return default script if file doesn't exist
    if network == "mainnet" {
        DEFAULT_MAINNET_SCRIPT.to_string()
    } else {
        DEFAULT_TESTNET_SCRIPT.to_string()
    }
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

// Add function to extract solana binary path
fn extract_solana_path(script_content: &str) -> Option<PathBuf> {
    if let Some(validator_line) = script_content.lines()
        .find(|line| line.contains("exec") && line.contains("solana-validator")) {
        if let Some(path) = validator_line.split("exec").nth(1) {
            if let Some(validator_path) = path.trim().split_whitespace().next() {
                return Some(PathBuf::from(validator_path).parent()?.join("solana"));
            }
        }
    }
    None
}

// Create and return the validator view layout
pub fn get_validator_view() -> LinearLayout {
    let initial_status = if is_validator_running() { 
        StyledString::styled("RUNNING", Style::from(Color::Dark(BaseColor::Green)))
    } else {
        StyledString::styled("STOPPED", Style::from(Color::Dark(BaseColor::Red)))
    };

    let dashboard = Panel::new(
        LinearLayout::vertical()
            .child(LinearLayout::horizontal()
                .child(TextView::new("Validator Status: "))
                .child(TextView::new(initial_status)
                    .with_name("status_text")))
            .child(LinearLayout::horizontal()
                .child(TextView::new("Catchup Status: "))
                .child(TextView::new(StyledString::styled("N/A", Style::from(Color::Dark(BaseColor::Yellow))))
                    .with_name("catchup_status_text")))
    )
    .title("Dashboard")
    .full_width()
    .fixed_height(5)
    .with_name("dashboard");

    // Add network selection radio group
    let mut radio_group = RadioGroup::new();
    let radio_button1 = radio_group.button("testnet".to_string(), "TESTNET");
    let radio_button2 = radio_group.button("mainnet".to_string(), "MAINNET");

    // Add callback for network selection
    radio_group.set_on_change(|s, network| {
        // Update current network
        if let Ok(mut current_network) = CURRENT_NETWORK.lock() {
            *current_network = network.to_string();
        }
        // Update script content
        s.call_on_name("script_content", |view: &mut TextArea| {
            view.set_content(get_script_content(network));
        });
    });

    let radio_layout = LinearLayout::horizontal()
        .child(TextView::new("Network: "))
        .child(radio_button1)
        .child(DummyView.fixed_width(2))
        .child(radio_button2)
        .with_name("network_layout");

    let text_area = TextArea::new()
        .content(get_script_content("testnet"))
        .disabled()
        .with_name("script_content")
        .min_height(10)
        .max_height(28);

    // Create button layout with space between buttons
    let button_layout = LinearLayout::horizontal()
        .child(Button::new("Edit Script", move |s| {
            // Check TextArea's current state
            let is_enabled = s.call_on_name("script_content", |view: &mut TextArea| {
                view.is_enabled()
            }).unwrap_or(false);
            
            // Log the current state
            update_logs(s, &format!("Current TextArea enabled state: {}", is_enabled));

            if !is_enabled {
                // TextArea is disabled, switch to edit mode
                update_logs(s, "Switching to edit mode...");
                s.call_on_name("script_content", |view: &mut TextArea| {
                    view.enable();  // Enable editing
                });
                s.call_on_name("edit_save_button", |view: &mut Button| {
                    view.set_label("Save Script");
                });
            } else {
                // TextArea is enabled, save and switch to view mode
                update_logs(s, "Saving script and switching to view mode...");
                save_script(s);
                s.call_on_name("script_content", |view: &mut TextArea| {
                    view.disable();  // Disable editing
                });
                s.call_on_name("edit_save_button", |view: &mut Button| {
                    view.set_label("Edit Script");
                });
            }
        }).with_name("edit_save_button"))
        .child(DummyView.fixed_width(4))
        .child(Button::new(if is_validator_running() { "Stop Validator" } else { "Start Validator" }, move |s| {
            if !is_validator_running() {
                // Get current network
                let network = CURRENT_NETWORK.lock()
                    .map(|network| network.clone())
                    .unwrap_or_else(|_| "testnet".to_string());

                // Get script path based on network
                if let Ok(exe_path) = env::current_exe() {
                    if let Some(exe_dir) = exe_path.parent() {
                        let script_name = if network == "mainnet" {
                            "validator-mainnet.sh"
                        } else {
                            "validator-testnet.sh"
                        };
                        let script_path = exe_dir.join(script_name);

                        // Check if validator script exists
                        if !script_path.exists() {
                            // Save the script using existing function
                            save_script(s);
                            update_logs(s, &format!("Created {} validator script from default content", network));
                        }
                    }
                }

                // Check if script is in edit mode
                let is_editing = s.call_on_name("script_content", |view: &mut TextArea| {
                    view.is_enabled()
                }).unwrap_or(false);

                if is_editing {
                    // Show dialog to remind user to save script first
                    s.add_layer(
                        Dialog::around(TextView::new("Please save the script before starting validator"))
                            .title("Save Required")
                            .button("Got it", |s| {
                                s.pop_layer();
                            })
                    );
                    return;
                }
            }
            toggle_run_stop(s);
        }).with_name("run_button"))
        .child(DummyView.fixed_width(4))
        .child(Button::new("Auto Check Status", move |s| {
            if IS_AUTO_CHECKING.load(Ordering::SeqCst) {
                // Stop auto checking
                IS_AUTO_CHECKING.store(false, Ordering::SeqCst);
                s.call_on_name("auto_check_button", |button: &mut Button| {
                    button.set_label("Auto Check Status");
                });
                update_logs(s, "Stopped auto checking status");
            } else {
                // Start auto checking
                IS_AUTO_CHECKING.store(true, Ordering::SeqCst);
                s.call_on_name("auto_check_button", |button: &mut Button| {
                    button.set_label("Stop Checking Status");
                });
                update_logs(s, "Started auto checking status");
                update_dashboard(s);
            }
        }).with_name("auto_check_button"))
        .child(DummyView.fixed_width(4))
        .child(Button::new("Check Validator Logs", move |s| {
            if is_validator_running() {
                if let Some(log_path) = extract_log_path(&get_script_content("testnet")) {
                    // Print log path to logs area
                    update_logs(s, &format!("Reading log file: {}", &log_path));
                    
                    // Create a thread to read logs
                    std::thread::spawn({
                        let log_path = log_path.clone();
                        let cb_sink = s.cb_sink().clone();
                        move || {
                            if let Ok(output) = Command::new("tail")
                                .args(["-n", "20", &log_path])  // Changed from 69 to 20
                                .output()
                            {
                                if let Ok(content) = String::from_utf8(output.stdout) {
                                    // Clean ANSI escape sequences before displaying
                                    let clean_content = clean_log_message(&content);
                                    // Send the content back to the main thread
                                    let _ = cb_sink.send(Box::new(move |s| {
                                        update_logs(s, "=== Validator Logs Start ===");
                                        update_logs(s, &clean_content);
                                        update_logs(s, "=== Validator Logs End ===");
                                    }));
                                }
                            }
                        }
                    });
                    
                    update_logs(s, "Log file refresh requested");
                }
            } else {
                update_logs(s, "Validator is not running");
            }
        }));

    let config_content = LinearLayout::vertical()
        .child(radio_layout)
        .child(DummyView.fixed_height(1))
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
fn save_script(s: &mut Cursive) {
    // Get current network from global variable
    let network = CURRENT_NETWORK.lock()
        .map(|network| Arc::new(network.clone()))
        .unwrap_or_else(|_| Arc::new("testnet".to_string()));

    // Get script content
    let content = s.call_on_name("script_content", |view: &mut TextArea| {
        view.get_content().trim_end().to_string()
    }).unwrap_or_default();

    // Get current executable path
    if let Ok(exe_path) = env::current_exe() {
        // Get the directory containing the executable
        if let Some(exe_dir) = exe_path.parent() {
            // Create script path based on network
            let script_name = if network.as_str() == "mainnet" {
                "validator-mainnet.sh"
            } else {
                "validator-testnet.sh"
            };
            let script_path = exe_dir.join(script_name);

            // Save the content to file
            match fs::write(&script_path, content) {
                Ok(_) => {
                    // Make the script executable (Unix-like systems only)
                    #[cfg(unix)]
                    {
                        use std::os::unix::fs::PermissionsExt;
                        if let Ok(metadata) = fs::metadata(&script_path) {
                            let mut perms = metadata.permissions();
                            perms.set_mode(0o755);
                            let _ = fs::set_permissions(&script_path, perms);
                        }
                    }
                    
                    IS_SCRIPT_MODIFIED.store(false, Ordering::SeqCst);
                    update_logs(s, &format!("Script {} saved successfully!", script_name));
                },
                Err(e) => {
                    update_logs(s, &format!("Failed to save script: {}", e));
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

#[allow(dead_code)]
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

#[allow(dead_code)]
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
    let is_running = is_validator_running();
    
    if !is_running {
        // Get current network
        let network = CURRENT_NETWORK.lock()
            .map(|network| network.clone())
            .unwrap_or_else(|_| "testnet".to_string());

        if let Ok(exe_path) = env::current_exe() {
            if let Some(exe_dir) = exe_path.parent() {
                // Create script path based on network
                let script_name = if network == "mainnet" {
                    "validator-mainnet.sh"
                } else {
                    "validator-testnet.sh"
                };
                let script_path = exe_dir.join(script_name);
                
                // execute the validator script
                match Command::new("bash")
                    .arg("-c")
                    .arg(format!("nohup {} &", script_path.display()))
                    .stdin(Stdio::null())
                    .stdout(Stdio::piped())
                    .stderr(Stdio::piped())
                    .spawn() {
                    Ok(mut child) => {
                        update_logs(siv, &format!("{} validator script started successfully! Please wait for status update...", 
                            if network == "mainnet" { "Mainnet" } else { "Testnet" }));
                        
                        // Get handles for stdout and stderr
                        let stdout = child.stdout.take().expect("Failed to capture stdout");
                        let stderr = child.stderr.take().expect("Failed to capture stderr");
                        
                        // Create new thread to handle stdout
                        let cb_sink = siv.cb_sink().clone();
                        std::thread::spawn(move || {
                            let reader = BufReader::new(stdout);
                            for line in reader.lines() {
                                if let Ok(line) = line {
                                    let _ = cb_sink.send(Box::new(move |s| {
                                        update_logs(s, &format!("Start command stdout: {}", line));
                                    }));
                                }
                            }
                        });
                        
                        // Create new thread to handle stderr
                        let cb_sink = siv.cb_sink().clone();
                        std::thread::spawn(move || {
                            let reader = BufReader::new(stderr);
                            for line in reader.lines() {
                                if let Ok(line) = line {
                                    let _ = cb_sink.send(Box::new(move |s| {
                                        update_logs(s, &format!("Start command stderr: {}", line));
                                    }));
                                }
                            }
                        });

                        // Create new thread to wait for process completion
                        std::thread::spawn(move || {
                            let _ = child.wait();  // Wait for process to finish without blocking output
                        });

                        // Async status update
                        let cb_sink = siv.cb_sink().clone();
                        std::thread::spawn(move || {
                            std::thread::sleep(std::time::Duration::from_secs(10));
                            let _ = cb_sink.send(Box::new(|s| {
                                // Start auto checking
                                IS_AUTO_CHECKING.store(true, Ordering::SeqCst);
                                s.call_on_name("auto_check_button", |button: &mut Button| {
                                    button.set_label("Stop Checking Status");
                                });
                                update_logs(s, "Started auto checking status");
                                update_dashboard(s);
                            }));
                        });

                        // Update button state to "Stop"
                        siv.call_on_name("run_button", |button: &mut Button| {
                            button.set_label("Stop Validator");
                        });
                    },
                    Err(e) => {
                        update_logs(siv, &format!("Failed to start {} validator: {}", 
                            if network == "mainnet" { "mainnet" } else { "testnet" }, e));
                    }
                }
            }
        }
    } else {
        // Get script content to get validator path and ledger path
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
                .stdout(Stdio::piped())
                .stderr(Stdio::piped())
                .output() {
                Ok(output) => {
                    let stdout = String::from_utf8_lossy(&output.stdout).into_owned();
                    let stderr = String::from_utf8_lossy(&output.stderr).into_owned();
                    
                    update_logs(siv, "Executing solana-validator exit command. Please wait for status update...");
                    if !stdout.is_empty() {
                        update_logs(siv, "Exit command stdout:");
                        update_logs(siv, &stdout);
                    }
                    if !stderr.is_empty() {
                        update_logs(siv, "Exit command stderr:");
                        update_logs(siv, &stderr);
                    }
                },
                Err(e) => {
                    update_logs(siv, &format!("Failed to execute validator exit command: {}", e));
                }
            }
        }
        
        let cb_sink = siv.cb_sink().clone();
        std::thread::spawn(move || {
            std::thread::sleep(std::time::Duration::from_secs(60));
            let _ = cb_sink.send(Box::new(|s| {
                // Stop auto checking
                IS_AUTO_CHECKING.store(false, Ordering::SeqCst);
                s.call_on_name("auto_check_button", |button: &mut Button| {
                    button.set_label("Auto Check Status");
                });
                update_dashboard(s);
            }));
        });
        
        // Update button state to "Start"
        siv.call_on_name("run_button", |button: &mut Button| {
            button.set_label("Start Validator");
        });
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

fn update_dashboard(siv: &mut Cursive) {
    let is_running = is_validator_running();
    
    // Add log output
    update_logs(siv, &format!("Checking validator status: {}", if is_running { "RUNNING" } else { "STOPPED" }));
    
    // Update the validator status text with color
    siv.call_on_name("status_text", |view: &mut TextView| {
        let styled_status = if is_running {
            StyledString::styled("RUNNING", Style::from(Color::Dark(BaseColor::Green)))
        } else {
            StyledString::styled("STOPPED", Style::from(Color::Dark(BaseColor::Red)))
        };
        view.set_content(styled_status);
    });

    // Immediately update catchup status based on validator state
    siv.call_on_name("catchup_status_text", |view: &mut TextView| {
        if !is_running {
            view.set_content(StyledString::styled("N/A", Style::from(Color::Dark(BaseColor::Yellow))));
        }
    });

    // Update button state based on validator state
    siv.call_on_name("run_button", |button: &mut Button| {
        if is_running {
            button.set_label("Stop Validator");
        } else {
            button.set_label("Start Validator");
        }
    });

    // Reset auto check status when validator is not running
    siv.call_on_name("auto_check_button", |button: &mut Button| {
        if !is_running {
            button.set_label("Auto Check Status");
            IS_AUTO_CHECKING.store(false, Ordering::SeqCst);
        }
    });

    // If validator is running, start periodic checks
    if is_running {
        // Get script content to extract solana path
        if let Some(script_content) = siv.call_on_name("script_content", |view: &mut TextArea| {
            view.get_content().to_string()
        }).as_deref() {
            if let Some(solana_path) = extract_solana_path(script_content) {
                // Create new thread for periodic checks
                let cb_sink = siv.cb_sink().clone();
                std::thread::spawn(move || {
                    // Initial wait for validator initialization
                    let _ = cb_sink.send(Box::new(|s| {
                        update_logs(s, "Waiting for a while before catchup check...");
                    }));
                    std::thread::sleep(std::time::Duration::from_secs(60));

                    while IS_AUTO_CHECKING.load(Ordering::SeqCst) {
                        // Check if validator is still running
                        if !is_validator_running() {
                            // Update validator status, catchup status, and button when validator stops
                            let _ = cb_sink.send(Box::new(|s| {
                                update_logs(s, "Detected validator has stopped running. Updating status...");
                                // Update validator status to STOPPED
                                s.call_on_name("status_text", |view: &mut TextView| {
                                    view.set_content(StyledString::styled(
                                        "STOPPED",
                                        Style::from(Color::Dark(BaseColor::Red))
                                    ));
                                });
                                
                                // Update catchup status to N/A
                                s.call_on_name("catchup_status_text", |view: &mut TextView| {
                                    view.set_content(StyledString::styled(
                                        "N/A",
                                        Style::from(Color::Dark(BaseColor::Yellow))
                                    ));
                                });

                                // Update button state to "Start Validator"
                                s.call_on_name("run_button", |button: &mut Button| {
                                    button.set_label("Start Validator");
                                });
                            }));
                            break;
                        }

                        // Check catchup status
                        let _ = cb_sink.send(Box::new(|s| {
                            update_logs(s, "Starting catchup status check...");
                        }));
                        
                        match Command::new("script")
                            .args([
                                "-f", 
                                "catchup.status", 
                                "-c", 
                                &format!("timeout 0.2s {} catchup --our-localhost", solana_path.display())
                            ])
                            .output() {
                            Ok(_) => {
                                // Try to read the status file
                                if let Ok(content) = fs::read_to_string("catchup.status") {
                                    // Check for error in content
                                    if content.contains("error") || content.contains("Error") {
                                        let _ = cb_sink.send(Box::new(move |s| {
                                            update_logs(s, "Catchup status: Error detected in output");
                                            s.call_on_name("catchup_status_text", |view: &mut TextView| {
                                                view.set_content(StyledString::styled(
                                                    "N/A",
                                                    Style::from(Color::Dark(BaseColor::Yellow))
                                                ));
                                            });
                                        }));
                                    } else if let Some(line) = content.lines()
                                        .find(|line| line.contains("slot(s) behind")) {
                                        let line_clone = line.to_string();
                                        
                                        if let Some(num_str) = line_clone.split_whitespace()
                                            .find(|&s| s.chars().all(|c| c.is_digit(10))) {
                                            if let Ok(slots_behind) = num_str.parse::<u64>() {
                                                let _ = cb_sink.send(Box::new(move |s| {
                                                    update_logs(s, &format!("Catchup status: {}", line_clone));
                                                    
                                                    // Update dashboard based on slots_behind value
                                                    s.call_on_name("catchup_status_text", |view: &mut TextView| {
                                                        if slots_behind > 0 {
                                                            view.set_content(StyledString::styled(
                                                                "BEHIND",
                                                                Style::from(Color::Dark(BaseColor::Red))
                                                            ));
                                                        } else {
                                                            view.set_content(StyledString::styled(
                                                                "CATCHUP",
                                                                Style::from(Color::Dark(BaseColor::Green))
                                                            ));
                                                        }
                                                    });
                                                }));
                                            }
                                        }
                                    } else {
                                        // No "slot(s) behind" found in content
                                        let _ = cb_sink.send(Box::new(move |s| {
                                            update_logs(s, "Catchup status: No slot information found");
                                            s.call_on_name("catchup_status_text", |view: &mut TextView| {
                                                view.set_content(StyledString::styled(
                                                    "N/A",
                                                    Style::from(Color::Dark(BaseColor::Yellow))
                                                ));
                                            });
                                        }));
                                    }
                                }

                                // Clean up the status file
                                let _ = std::fs::remove_file("catchup.status");
                            },
                            Err(e) => {
                                let err_msg = e.to_string();
                                let _ = cb_sink.send(Box::new(move |s| {
                                    update_logs(s, &format!("Failed to execute catchup command: {}", err_msg));
                                }));
                            }
                        }

                        // Wait for 600 seconds before next check
                        let _ = cb_sink.send(Box::new(|s| {
                            update_logs(s, "Waiting 600 seconds before next catchup check...");
                        }));
                        std::thread::sleep(std::time::Duration::from_secs(600));
                    }

                    // When auto-checking stops, reset catchup status to N/A
                    let _ = cb_sink.send(Box::new(|s| {
                        s.call_on_name("catchup_status_text", |view: &mut TextView| {
                            view.set_content(StyledString::styled("N/A", Style::from(Color::Dark(BaseColor::Yellow))));
                        });
                    }));

                    // Update button state when auto-checking stops
                    let _ = cb_sink.send(Box::new(|s| {
                        s.call_on_name("auto_check_button", |button: &mut Button| {
                            button.set_label("Auto Check Status");
                        });
                        IS_AUTO_CHECKING.store(false, Ordering::SeqCst);
                        update_logs(s, "Stopped auto checking status");
                    }));
                });
            }
        }
    }
}