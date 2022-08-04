// Polyfill instantiateStreaming for browsers missing it
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

// Create promise to handle Worker calls whilst
// module is still initialising
let wasmResolve;
let wasmReady = new Promise((resolve) => {
    wasmResolve = resolve;
})

// Handle incoming messages
self.addEventListener('message', function(event) {

    const { eventType, eventData, eventId } = event.data;

    if (eventType === "INITIALISE") {
        // initialize the Go WASM glue
        self.importScripts('wasm_exec.js');
        const go = new self.Go();
        const result = WebAssembly.instantiateStreaming(fetch(eventData), go.importObject)
            .then(instantiatedModule => {
                go.run(instantiatedModule.instance)
                const wasmExports = instantiatedModule.instance.exports;

                // Resolve our exports for when the messages
                // to execute functions come through
                wasmResolve(wasmExports);

                // Send back initialised message to main thread
                self.postMessage({
                    eventType: "INITIALISED",
                    eventData: Object.keys(wasmExports)
                });

            });
    } else if (eventType === "CALL") {
        wasmReady
            .then((wasmInstance) => {
                const method = wasmInstance[eventData.method];
                const result = method.apply(null, eventData.arguments);
                self.postMessage({
                    eventType: "RESULT",
                    eventData: result,
                    eventId: eventId
                });
            })
            .catch((error) => {
                self.postMessage({
                    eventType: "ERROR",
                    eventData: "An error occured executing WASM instance function: " + error.toString(),
                    eventId: eventId
                });
            })
    }

}, false);