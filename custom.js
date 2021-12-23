const helper = {
    customize: (selector, property, newValue) => {
        let el = document.querySelector(selector);
        if (el?.style) el.style[property] = newValue;
    }
}


//students
helper.customize("#page-content > div > div.l-content-main.proto > div:nth-child(10)", "display", "none");

//lectutets
helper.customize("#page-content > div > div.l-content-main.proto > table.no_padding", "display", "none");
helper.customize("#page-content > div > div.l-content-main.proto > div.box2.t_gray_light.t_small.margin_top_x", "display", "none");
helper.customize("#page-content > div > div.l-content-main.proto > div.box2.t_gray_light.t_small.margin_top", "display", "none");

if (!document.querySelector(".l-content-main"))
    document.body.classList.add(".l-content-main")
