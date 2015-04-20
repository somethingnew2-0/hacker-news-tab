var page = require('webpage').create(),
    system = require('system'),
    address, size;

if (system.args.length < 2 || system.args.length > 4) {
    console.error('Usage: rasterize.js URL [paperwidth*paperheight|paperformat] [zoom]');
    console.error('  image (png/jpg output) examples: "1920px" entire page, window width 1920px');
    console.error('                                   "800px*600px" window, clipped to 800x600');
    phantom.exit(1);
} else {
    address = system.args[1];
    page.viewportSize = { width: 600, height: 600 };
    if (system.args.length > 2 && system.args[2].substr(-2) === "px") {
        size = system.args[2].split('*');
        if (size.length === 2) {
            pageWidth = parseInt(size[0], 10);
            pageHeight = parseInt(size[1], 10);
            page.viewportSize = { width: pageWidth, height: pageHeight };
            page.clipRect = { top: 0, left: 0, width: pageWidth, height: pageHeight };
        } else {
            console.error("size:", system.args[2]);
            pageWidth = parseInt(system.args[2], 10);
            pageHeight = parseInt(pageWidth * 3/4, 10); // it's as good an assumption as any
            console.error("pageHeight:",pageHeight);
            page.viewportSize = { width: pageWidth, height: pageHeight };
        }
    }
    if (system.args.length > 3) {
        page.zoomFactor = system.args[3];
    }
    page.open(address, function (status) {
        if (status !== 'success') {
            console.error('Unable to load the address!');
            phantom.exit(1);
        } else {
            window.setTimeout(function () {
                var base64 = page.renderBase64('PNG');
                console.log(base64);
                phantom.exit();
            }, 500);
        }
    });
}
