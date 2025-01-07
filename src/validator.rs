use cursive::views::{LinearLayout, Panel, TextView};
use cursive::traits::*;

// Create and return the validator view layout
pub fn get_validator_view() -> LinearLayout {
    // Create three sections
    let dashboard = Panel::new(TextView::new("Validator Dashboard"))
        .title("Dashboard")
        .full_width()
        .fixed_height(5);

    let config = Panel::new(TextView::new("Validator Config"))
        .title("Config")
        .with_name("right_panel")
        .full_width()
        .full_height();

    let logs = Panel::new(TextView::new("Validator Logs"))
        .title("Logs")
        .full_width()
        .fixed_height(8);

    // Combine sections vertically
    let layout = LinearLayout::vertical()
        .child(dashboard)
        .child(config)
        .child(logs);
    
    // Return the layout
    layout
}