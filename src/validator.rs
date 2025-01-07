use cursive::views::{LinearLayout, Panel, TextView, TextArea, Button};
use cursive::traits::*;
use std::sync::Arc;
use std::sync::atomic::{AtomicBool, Ordering};
use std::env;
use std::fs;
use std::path::PathBuf;
use std::io::Write;

// Global state for run/stop button
static IS_RUNNING: AtomicBool = AtomicBool::new(false);

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
"#;

// Create and return the validator view layout
pub fn get_validator_view() -> LinearLayout {
    // Create three sections
    let dashboard = Panel::new(TextView::new("Validator Dashboard"))
        .title("Dashboard")
        .full_width()
        .fixed_height(5);

    // Create config section with TextArea and buttons
    let text_area = TextArea::new()
        .content(DEFAULT_SCRIPT)
        .with_name("script_content")
        .min_height(10);

    let button_layout = LinearLayout::horizontal()
        .child(Button::new("Save", |s| {
            save_script(s);
        }))
        .child(Button::new("Run", move |s| {
            toggle_run_stop(s);
        }).with_name("run_button"));

    let config = Panel::new(
        LinearLayout::vertical()
            .child(text_area)
            .child(button_layout)
    )
    .title("Config")
    .full_width()
    .full_height();

    let logs = Panel::new(TextView::new("Validator Logs"))
        .title("Logs")
        .with_name("log_view")
        .full_width()
        .fixed_height(8);

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
                    update_logs(siv, "Script saved successfully!");
                },
                Err(e) => {
                    // Update log to show error
                    update_logs(siv, &format!("Failed to save script: {}", e));
                }
            }
        }
    }
}

// Function to toggle between Run and Stop
fn toggle_run_stop(siv: &mut cursive::Cursive) {
    let is_running = IS_RUNNING.load(Ordering::SeqCst);
    
    if !is_running {
        // Start the script
        // TODO: Execute the saved script
        // Will implement when you provide the command
        
        siv.call_on_name("run_button", |button: &mut Button| {
            button.set_label("Stop");
        });
        IS_RUNNING.store(true, Ordering::SeqCst);
    } else {
        // Stop the script
        // TODO: Execute stop command
        // Will implement when you provide the command
        
        siv.call_on_name("run_button", |button: &mut Button| {
            button.set_label("Run");
        });
        IS_RUNNING.store(false, Ordering::SeqCst);
    }
}

// Function to update log view
fn update_logs(siv: &mut cursive::Cursive, message: &str) {
    siv.call_on_name("log_view", |view: &mut TextView| {
        view.append(format!("{}\n", message));
    });
}