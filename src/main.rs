mod validator;

use cursive::Cursive;
use cursive::theme::{Theme, BaseColor, Color, PaletteColor, ColorStyle};
use cursive::views::{LinearLayout, SelectView, Panel, TextView};
use cursive::traits::*;
use cursive::event::Event;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::time::{SystemTime, Duration};
use lazy_static::lazy_static;
use cursive::utils::markup::StyledString;

lazy_static! {
    static ref Q_COUNT: AtomicUsize = AtomicUsize::new(0);
    static ref LAST_Q_TIME: std::sync::Mutex<SystemTime> = std::sync::Mutex::new(SystemTime::now());
}

// Handle menu item selection
fn menu_selected(siv: &mut Cursive, item: &str) {
    match item {
        "validator" => {
            // Replace right panel content with validator view
            siv.call_on_name("right_sections", |view: &mut LinearLayout| {
                *view = validator::get_validator_view();
            });
        },
        "quit_info" => {
            // Do nothing for quit info item
        },
        _ => {
            siv.call_on_name("right_panel", |view: &mut Panel<TextView>| {
                view.get_inner_mut().set_content(format!("Selected: {}", item));
            });
        }
    }
}

fn main() {
    // Initialize the cursive interface
    let mut siv = cursive::default();
    
    // Disable Ctrl-c
    siv.clear_global_callbacks(Event::CtrlChar('c'));
    
    // Add 'q' key handler for quitting
    siv.add_global_callback('q', |s| {
        let now = SystemTime::now();
        let mut last_time = LAST_Q_TIME.lock().unwrap();
        
        // If more than 2 seconds have passed, reset the counter
        if now.duration_since(*last_time).unwrap_or(Duration::from_secs(0)) > Duration::from_secs(2) {
            Q_COUNT.store(0, Ordering::SeqCst);
        }
        
        // Update last press time
        *last_time = now;
        
        // Increment counter
        let count = Q_COUNT.fetch_add(1, Ordering::SeqCst) + 1;
        
        if count >= 4 {
            s.quit();
        }
    });
    
    // Set up the theme with unified black background
    let mut theme = Theme::default();
    theme.palette[PaletteColor::Background] = Color::Dark(BaseColor::Black);
    theme.palette[PaletteColor::View] = Color::Dark(BaseColor::Black);
    theme.palette[PaletteColor::Primary] = Color::Light(BaseColor::White);
    theme.palette[PaletteColor::Secondary] = Color::Light(BaseColor::White);
    theme.palette[PaletteColor::Shadow] = Color::Dark(BaseColor::Black);
    siv.set_theme(theme);

    // Create menu view with selection handler
    let mut menu = SelectView::new()
        .on_submit(menu_selected);
    
    // Add menu items
    menu.add_item("X1 Validator", "validator");
    menu.add_item("", "");  // Add empty item as spacer
    menu.add_item(
        StyledString::styled(
            "QUIT (Press 'q' 4 times)",
            ColorStyle::new(Color::Light(BaseColor::Red), Color::Dark(BaseColor::Black))
        ),
        "quit_info"
    );
    menu.add_item("", "");  // Add empty item as spacer
    menu.add_item(
        StyledString::styled(
            "by xen_artist",
            ColorStyle::new(Color::Dark(BaseColor::Green), Color::Dark(BaseColor::Black))
        ),
        "author_info"
    );

    // Set default selection to X1 Validator
    menu.set_selection(0);
    
    // Create left panel with title
    let left_panel = Panel::new(menu)
        .title("xoon")
        .min_width(20)
        .full_height();

    // Create initial right panel with validator view
    let right_sections = validator::get_validator_view()
        .with_name("right_sections")
        .full_width()
        .full_height();

    // Create main layout with left menu and right content panel
    let layout = LinearLayout::horizontal()
        .child(left_panel)
        .child(right_sections)
        .full_width()
        .full_height();

    // Add the layout to the screen
    siv.add_layer(layout);
    
    // Start the event loop
    siv.run();
}