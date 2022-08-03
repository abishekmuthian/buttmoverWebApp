const button = document.querySelector(".connect");
button.addEventListener("click", async function () {
    if ("serial" in navigator) {
        // The Web Serial API is supported.
        // Filter on devices with the Arduino Uno USB Vendor/Product IDs.
        const filters = [{ usbVendorId: 0x1a86, usbProductId: 0x7523 }];

        // Prompt user to select an Arduino Uno device.
        const port = await navigator.serial.requestPort({ filters });

        const { usbProductId, usbVendorId } = port.getInfo();

        // Wait for the serial port to open.
        await port.open({ baudRate: 115200 });

        // Listen to data coming from the serial device.
        const textDecoder = new TextDecoderStream();
        const readableStreamClosed = port.readable.pipeTo(textDecoder.writable);
        const reader = textDecoder.readable.getReader();



        // Listen to data coming from the serial device.
        while (port.readable) {
            let timer25 = 1500;
            let timer5 = 300;
            let timer30 = 1800;
            let mover25Id;
            let break5Id;
            let break30Id;
            let movercompleted = false;
            let movercount = 0;
            try {
                while (true) {
                    const { value, done } = await reader.read();
                    if (done) {
                        // Allow the serial port to be closed later.
                        reader.releaseLock();
                        break;
                    }
                    if (value) {
                        console.log(value);
                        if (value.includes("Button Pressed")){
                            window.clearInterval(break5Id);
                            timer5 = 300;
                            window.clearInterval(break30Id);
                            timer30 = 1800;
                            const p = tick25();
                            const message = await p
                            console.log(message)
                            mover25Id = window.setInterval(async function() {
                                // this callback will execute every second until we call
                                // the clearInterval method
                                timer25--;
                                console.log(timer25);
                                if (timer25 === 0) {
                                    // TODO: go ahead and really log the dude out
                                    const stoptick = stoptick25();
                                    const stoptickmessage = await stoptick
                                    console.log(stoptickmessage)


                                    const p = mover25();
                                    const message = await p
                                    console.log(message)
                                    movercount++;
                                    movercompleted = true;
                                    window.clearInterval(mover25Id);
                                    timer25 = 1500;
                                    //alert(message)
                                }
                            }, 1000);

                        }
                        if (value.includes("Button Released")){
                            window.clearInterval(mover25Id);
                            timer25 = 1500;
                            if (movercount !== 0 && movercount < 4 && movercompleted) {
                                break5Id = window.setInterval(async function() {
                                    // this callback will execute every second until we call
                                    // the clearInterval method
                                    timer5--;
                                    console.log(timer5);
                                    if (timer5 === 0) {
                                        // TODO: go ahead and really log the dude out
                                        const p = break5();
                                        const message = await p;
                                        console.log(message);
                                        window.clearInterval(break5Id);
                                        timer5 = 300;
                                        //alert(message)
                                    }
                                }, 1000);
                            } else if (movercount != 0 && movercount >= 4 && movercompleted) {
                                break30Id = window.setInterval(async function() {
                                    // this callback will execute every second until we call
                                    // the clearInterval method
                                    timer30--;
                                    console.log(timer30);
                                    if (timer5 === 0) {
                                        // TODO: go ahead and really log the dude out
                                        const p = break30();
                                        const message = await p;
                                        console.log(message);
                                        window.clearInterval(break30Id);
                                        timer30 = 1800;
                                        //alert(message)
                                    }
                                }, 1000);
                            }else{
                                const p = break0();
                                const message = await p;
                                console.log(message);
                                //alert(message)
                            }
                        }
                    }
                }
            } catch (error) {
                // TODO: Handle non-fatal read error.
                console.log(error);
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
        });
    } else {
        // Browser not supported
        console.log("Browser not supported");
    }
});


function connectButtTrigger(){

}

function buttonPressed(){

}

function buttonReleased(){

}