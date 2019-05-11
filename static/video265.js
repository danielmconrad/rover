const videoURL = `ws://${location.host}/video`;
const videoContainer = document.getElementById('video-container');
const videoCanvas = document.getElementById('video-canvas');

console.log("Playing with libde265", libde265.de265_get_version());

const player = new libde265.StreamPlayer(videoCanvas);

player.set_status_callback(function(msg, fps) {
    player.disable_filters(true);

    switch (msg) {
    case "loading":
        console.log("Loading movie...");
        break;
    case "initializing":
        console.log("Initializing...");
        break;
    case "playing":
        console.log("Playing...");
        break;
    case "stopped":
        console.log("Stopped...");
        break;
    case "fps":
        console.log(Number(fps).toFixed(2) + " fps");
        break;
    default:
        console.log(msg);
    }
});

player.playback(videoURL);
