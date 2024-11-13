export var Stop = false;
var frameCount = 0;
var fpsInterval, startTime, now, then, elapsed;

export function StartAnimating(fps) {
    fpsInterval = 1000 / fps;
    then = Date.now();
    startTime = then;
    animate();
}


function animate() {

    // stop
    if (Stop) {
        return
    }
    let fpsDiv = document.getElementById("FPS-counter")
    if (fpsDiv == undefined) {
        return
    }
    // request another frame
    requestAnimationFrame(animate)

    // calc elapsed time since last loop

    now = Date.now();
    elapsed = now - then;

    // if enough time has elapsed, draw the next frame
    if (elapsed > fpsInterval) {

        // Get ready for next frame by setting then=now, but...
        // Also, adjust for fpsInterval not being multiple of 16.67
        then = now - (elapsed % fpsInterval);

        // TESTING...Report #seconds since start and achieved fps.
        var sinceStart = now - startTime;
        var currentFps = Math.round(1000 / (sinceStart / ++frameCount) * 100) / 100;
        fpsDiv.innerHTML = "Elapsed time= " + Math.round(sinceStart / 1000 * 100) / 100 + " secs @ " + currentFps + " fps."

    }
}