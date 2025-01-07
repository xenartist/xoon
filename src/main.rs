use cursive::Cursive;
use cursive::theme::{Theme, BaseColor, Color, PaletteColor};
use cursive::views::{LinearLayout, SelectView, Panel, TextView};
use cursive::traits::*;

// Handle menu item selection
fn menu_selected(siv: &mut Cursive, item: &str) {
    // Get the right panel by name and update its content
    siv.call_on_name("right_panel", |view: &mut Panel<TextView>| {
        view.get_inner_mut().set_content(format!("Selected: {}", item));
    });
}

fn main() {
    // Initialize cursive
    let mut siv = cursive::default();

    // Set theme
    let mut theme = Theme::default();
    theme.palette[PaletteColor::Background] = Color::Dark(BaseColor::Black);
    siv.set_theme(theme);

    // Create menu view with selection handler
    let mut menu = SelectView::new()
        .on_submit(menu_selected);
    
    // Add menu items
    menu.add_item("X1 Validator", "menu1");

    // Create right panel with default content
    let right_panel = Panel::new(TextView::new("Select a menu item"))
        .with_name("right_panel")
        .full_width()   // Make right panel use all available width
        .full_height(); // Make right panel use all available height

    // Create main layout with left menu and right content panel
    let layout = LinearLayout::horizontal()
        .child(Panel::new(menu).min_width(20).full_height())  // Left panel with minimum width
        .child(right_panel)
        .full_width()   // Make layout use all available width
        .full_height(); // Make layout use all available height

    // Add the layout to the screen
    siv.add_layer(layout);
    
    // Start the event loop
    siv.run();
}