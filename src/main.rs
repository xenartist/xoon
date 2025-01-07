use cursive::Cursive;
use cursive::theme::{Theme, BaseColor, Color, PaletteColor};
use cursive::views::{LinearLayout, SelectView, Panel, TextView};
use cursive::traits::*;

// Handle menu item selection
fn menu_selected(siv: &mut Cursive, item: &str) {
    siv.call_on_name("right_panel", |view: &mut Panel<TextView>| {
        view.get_inner_mut().set_content(format!("Selected: {}", item));
    });
}

fn main() {
    // Initialize the cursive interface
    let mut siv = cursive::default();

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
    menu.add_item("X1 Validator", "menu1");

    let left_panel = Panel::new(menu)
        .title("xoon")
        .min_width(20)
        .full_height();

    // Create three sections for right panel
    let dashboard = Panel::new(TextView::new("Dashboard"))
        .title("Dashboard")
        .full_width()
        .fixed_height(5);

    let config = Panel::new(TextView::new("Config"))
        .title("Config")
        .with_name("right_panel")
        .full_width()
        .full_height();

    let logs = Panel::new(TextView::new("Log"))
        .title("Logs")
        .full_width()
        .fixed_height(8);

    // Combine right sections vertically
    let right_sections = LinearLayout::vertical()
        .child(dashboard)
        .child(config)
        .child(logs)
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