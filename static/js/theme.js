
document.addEventListener('DOMContentLoaded', () => {
    const themeController = document.querySelector('html');
    const savedTheme = localStorage.getItem('theme') || 'light';
    themeController.setAttribute('data-theme', savedTheme);

    const themeOptions = document.querySelectorAll('input[name="theme-dropdown"]');
    themeOptions.forEach(option => {
        if (option.value === savedTheme) {
            option.checked = true;
        }
        option.addEventListener('change', (e) => {
            const newTheme = e.target.value;
            themeController.setAttribute('data-theme', newTheme);
            localStorage.setItem('theme', newTheme);
        });
    });
});
