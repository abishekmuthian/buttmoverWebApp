// Avoid CORS
// document.domain = "localhost";
var parentWindow = window.parent;

var gameOver = false;

// Handle worker
if (typeof (worker) == "undefined") {
    worker = new Worker("./js/timer.worker.js");
    worker.addEventListener('message', this.handleMessage)
}

// Handling stream chunks
class LineBreakTransformer {
    constructor() {
        // A container for holding stream data until a new line.
        this.chunks = "";
    }

    transform(chunk, controller) {
        // Append new chunks to existing chunks.
        this.chunks += chunk;
        // For each line breaks in chunks, send the parsed lines out.
        const lines = this.chunks.split("\r\n");
        this.chunks = lines.pop();
        lines.forEach((line) => controller.enqueue(line));
    }

    flush(controller) {
        // When the stream is closed, flush any remaining chunks out.
        controller.enqueue(this.chunks);
    }
}

// Helper functions
function callConnectButtTrigger() {
    parentWindow.connectSerial().then(async navigator => {
        if ("serial" in navigator) {
            console.log("Inside serial connection")
            // The Web Serial API is supported.
            // Filter on devices with the Arduino Uno USB Vendor/Product IDs.
            const filters = [{usbVendorId: 0x1a86, usbProductId: 0x7523}];

            try {
                // Prompt user to select an Arduino Uno device.
                const port = await navigator.serial.requestPort({filters});
                const {usbProductId, usbVendorId} = port.getInfo();

                // Wait for the serial port to open.
                await port.open({baudRate: 115200});

                // const [appReadable, devReadable] = port.readable.tee();
                // console.log(devReadable)
                // You may want to update UI with incoming data from appReadable
                // and log incoming data in JS console for inspection from devReadable.

                // Listen to data coming from the serial device.
                const textDecoder = new TextDecoderStream();
                const readableStreamClosed = port.readable.pipeTo(textDecoder.writable);
                // const reader = textDecoder.readable.getReader();
                const reader = textDecoder.readable
                    .pipeThrough(new TransformStream(new LineBreakTransformer()))
                    .getReader();

                // Listen to data coming from the serial device.
                while (true) {
                    setConnectionStatus(true);
                    try {
                        while (true) {
                            const {value, done} = await reader.read();
                            if (done) {
                                // Allow the serial port to be closed later.
                                reader.releaseLock();
                                break;
                            }
                            if (value) {
                                console.log("Serial: ", value);
                                if (value.includes("Button Pressed")) {
                                    buttonPressed()
                                }
                                if (value.includes("Button Released")) {
                                    buttonReleased()
                                }
                            }
                        }
                    } catch (error) {
                        // TODO: Handle non-fatal read error.
                        console.log("Port issue: ", error);
                    }
                }

                navigator.serial.addEventListener("connect", (event) => {
                    // TODO: Automatically open event.target or warn user a port is available.
                    console.log("Butt Mover connected");
                });

                navigator.serial.addEventListener("disconnect", (event) => {
                    // TODO: Remove |event.target| from the UI.
                    // If the serial port was opened, a stream error would be observed as well.
                    console.log("Butt Mover disconnected");
                    setErrorStatus("Butt Trigger has been lost");
                });
            } catch (error) {
                console.log("Error in selecting the port, ", error);
                setErrorStatus("Error in selecting the port");
            }

        } else {
            // Browser not supported
            console.log("Browser not supported");
            setErrorStatus("Browser not supported");
        }
    });
}

async function handleMessage(message) {
    switch (message.data.event) {
        case 'setTick':
            const settickmessage = setTick(message.data.minute, message.data.second, message.data.moveYourAss, message.data.bringYourAssBack, message.data.timerLabel)
            console.log(settickmessage)
            break;
        case 'mover25Completed':
            const stoptick = stopTick25();
            const stoptickmessage = await stoptick;
            console.log(stoptickmessage);

            const m25 = mover25();
            const m25Message = await m25;
            console.log(m25Message);
            break;
        case  'tick30':
            const t30 = tick30();
            const t30Message = await t30;
            console.log(t30Message);
            break;
        case 'tick5':
            const t5 = tick5();
            const t5Message = await t5;
            console.log(t5Message);
            break;
        case 'break5':
            const b5 = break5();
            const b5Message = await b5;
            console.log(b5Message);
            break;
        case 'break30':
            const b30 = break30();
            const b30Message = await b30;
            console.log(b30Message);
            break;
        case 'break0':
            const stoptick25 = stopTick25();
            const stoptick25message = await stoptick25;
            console.log(stoptick25message);

            const b0 = break0();
            const b0Message = await b0;
            console.log(b0Message);
            break;
        case 'scoreUpdate':
            const ss = setScore(message.data.breaks, message.data.moves, message.data.lifes, message.data.newg);
            const ssmessage = await ss;
            console.log(ssmessage)
            break;
        case 'metricUpdate':
            console.log("metricUpdate");
            if (message.data.task === true) {

                const tasksContainer = parent.document.querySelector(".tasks");

                if (tasksContainer !== undefined) {
                    let tasks = parseInt(tasksContainer.innerText);
                    var newTasks = tasks + 1;

                    // Set breaks in cookie
                    parent.document.cookie = "tasks=" + newTasks;

                    tasksContainer.innerText = newTasks;

                    let taskHourContainer = parent.document.querySelector(".task_hours");
                    let taskHourMinutes = newTasks * 25;
                    let taskHour = minuteToHour(taskHourMinutes);
                    console.log("Task Hour: ", taskHour)
                    taskHourContainer.innerText = taskHour;
                }

                // Get authenticity token
                let token = parentWindow.authenticityToken();

                // Perform a post to the specified url (href of link)
                let url = "/users/task/increment";
                let data = "authenticity_token=" + token;

                parentWindow.DOM.Post(url, data, function (request) {
                    // Do something with response

                }, function (request) {
                    // Respond to error
                    console.log("error", request);
                });
            } else if (message.data.break === true) {

                let breaksContainer = parent.document.querySelector(".breaks");

                if (breaksContainer !== undefined) {
                    let breaks = parseInt(breaksContainer.innerText);
                    var newBreaks = breaks + 1;


                    // Set breaks in cookie
                    parent.document.cookie = "breaks=" + newBreaks;

                    breaksContainer.innerText = newBreaks;

                    let breakHourContainer = parent.document.querySelector(".break_hours");

                    let breakHourMinutes;

                    if (newBreaks >= 5) {
                        let thirtyMinuteBreaksDifference = newBreaks % 5
                        let thirtyMinuteBreaks = (newBreaks - thirtyMinuteBreaksDifference) / 5;
                        let fiveMinuteBreaks = newBreaks - thirtyMinuteBreaks;
                        breakHourMinutes = (fiveMinuteBreaks * 5) + (thirtyMinuteBreaks * 30);
                    } else {
                        breakHourMinutes = newBreaks * 5;
                    }
                    let breakHour = minuteToHour(breakHourMinutes);
                    console.log("Break Hour: ", breakHour)
                    breakHourContainer.innerText = breakHour;
                }

                // Get authenticity token
                let token = parentWindow.authenticityToken();

                // Perform a post to the specified url (href of link)
                let url = "/users/break/increment";
                let data = "authenticity_token=" + token;

                parentWindow.DOM.Post(url, data, function (request) {
                    // Do something with response

                }, function (request) {
                    // Respond to error
                    console.log("error", request);
                });
            }

    }
}

function minuteToHour(totalMinutes) {
    var hour = Math.floor(totalMinutes / 60);
    var minutes = totalMinutes % 60;

    if (hour === 0) {
        return minutes + "m";
    }

    if (minutes === 0) {
        return hour + "h"
    }

    return hour + "h" + minutes + "m"
}

async function buttonPressed() {
    console.log("Button Pressed");
    clearBreak();
    clearTask();

    if (gameOver) {
        reset();
    }

    /*    const stoptick30 = stopTick30();
        const stoptick30message = await stoptick30;
        console.log(stoptick30message);

        const stoptick5 = stopTick5();
        const stoptick5message = await stoptick5;
        console.log(stoptick5message);*/

    setClock();

    worker.postMessage({event: 'task'});

    setButtonTrigger(true);

    /*    const p = tick25();
        const message = await p
        console.log(message)*/

}

async function buttonReleased() {
    console.log("Button Released");
    if (gameOver) {
        console.log("Game Over flag detected, So doing nothing");
        return;
    }
    clearTask();
    clearBreak();

    worker.postMessage({event: 'break', mode: 'buttonrelease'});

    setButtonTrigger(false);
}

async function clearTask() {
    console.log("Clear Task");
    worker.postMessage({event: 'cleartask'});
}

async function clearBreak() {
    console.log("Clear Break");
    worker.postMessage({event: 'clearbreak'})
}

async function reset() {
    console.log("Reset");
    worker.postMessage({event: 'reset'});
}

async function setGameOverFlag(arg) {
    gameOver = arg;
}


// Polyfill
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("./wasm/main.wasm"), go.importObject).then(result => {
    document.getElementById('loading').remove();
    go.run(result.instance);

    // Set and detect theme
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        // dark mode
        console.log("Dark mode");
        setColorMode("dark");
    }

    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
        // light mode
        console.log("Light mode");
        setColorMode("light");
    }

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
        const newColorScheme = event.matches ? "dark" : "light";
        console.log(newColorScheme + " " + "mode");
        setColorMode(newColorScheme);
    });
});

// Set Game color mode
function setLightColorMode() {
    setColorMode("light");
}

function setDarkColorMode() {
    setColorMode("dark");
}






