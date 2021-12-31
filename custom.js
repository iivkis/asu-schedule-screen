const helper = {
    customize: (selector, property, newValue) => {
        let el = document.querySelector(selector);
        if (el?.style) el.style[property] = newValue;
    },
    isLocation: (pathRegexp) => {
        return pathRegexp.test(window.location.pathname)
    }
}


if (!document.querySelector(".l-content-main"))
    document.body.classList.add(".l-content-main")


//скрываем лишние элементы
if (helper.isLocation(/students/ig)) {
    helper.customize("#page-content > div > div.l-content-main.proto > div:nth-child(11)", "display", "none");
    helper.customize("#page-content > div > div.l-content-main.proto > div.box2.t_gray_light.t_small.margin_top_x", "display", "none")
}

if (helper.isLocation(/lecturers/ig)) {
    helper.customize("#page-content > div > div.l-content-main.proto > table", "display", "none")
    helper.customize("#page-content > div > div.l-content-main.proto > div.box2.t_gray_light.t_small.margin_top_x", "display", "none")
}


