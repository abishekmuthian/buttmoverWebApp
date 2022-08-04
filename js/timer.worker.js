let timer25 = 1500;
let timer5 = 300;
let timer30 = 1800;
let minute = 0;
let second = 0;
let task25Id;
let break5Id;
let break30Id;
let taskcompleted = false;
let taskcount = 0;
let breakcount = 0;
let lifecount = 0;
let buttontrigger = false;
let newgame = true;

self.onmessage = function (msg) {
    switch (msg.data.mode) {
        case 'buttonpress':
            buttontrigger = true;
            break;
        case 'buttonrelease':
            buttontrigger = false;
            break;
        default:
            break;
    }
    switch (msg.data.event) {
        case 'clearbreak':
            handleClearBreak();
            break;
        case 'cleartask':
            handleClearTask();
            break;
        case 'task':
            handleTask();
            break;
        case 'break':
            handleBreak();
            break;
        case 'reset':
            handleReset();
        default:
            break;
    }
}

function handleReset() {
    newgame = true;


    taskcount = 0;
    breakcount = 0;
    lifecount = 100;


    handleClearTask();
    handleClearBreak();
}

function handleClearBreak() {
    console.log("Inside handleClearBreak");
    clearInterval(break5Id);
    timer5 = 300;

    clearInterval(break30Id);
    timer30 = 1800;

    minute = 0;
    second = 0;

    if (newgame) {
        lifecount = 100;
    } else {
        lifecount = (100 * breakcount) / taskcount;
    }

    self.postMessage({
        event: 'scoreUpdate',
        breaks: breakcount,
        moves: taskcount,
        lifes: lifecount,
        newg: newgame,
    });

}

function handleClearTask() {
    console.log("Inside handleClearMover");

    clearInterval(task25Id);
    console.log("Cleared mover25 in handleClearMover");
    timer25 = 1500;

    minute = 0;
    second = 0;

    if (newgame) {
        lifecount = 100;
    } else {
        lifecount = (100 * breakcount) / taskcount;
    }

    self.postMessage({
        event: 'scoreUpdate',
        breaks: breakcount,
        moves: taskcount,
        lifes: lifecount,
        newg: newgame,
    });
}

function handleTask() {

    task25Id = setInterval(async function () {
        // this callback will execute every second until we call
        // the clearInterval method
        taskcompleted = false;
        timer25--;
        console.log(timer25);

        second++;

        if (second > 59) {
            second = 0
            minute++
        }

        self.postMessage({
            event: 'setTick',
            minute: minute,
            second: second,
            moveYourAss: false,
            bringYourAssBack: false,
            timerLabel: 'Task',
        });

        if (timer25 === 0) {
            console.log("Increasing mover count")
            taskcount++;
            lifecount = (100 * breakcount) / taskcount;

            if (taskcount === 1) {
                newgame = true;
                lifecount = 100;
                console.log("First break")

                self.postMessage({
                    event: 'scoreUpdate',
                    breaks: breakcount,
                    moves: taskcount,
                    lifes: lifecount,
                    newg: newgame,
                });
                lifecount = 100;
            } else {
                newgame = false;

                self.postMessage({
                    event: 'scoreUpdate',
                    breaks: breakcount,
                    moves: taskcount,
                    lifes: lifecount,
                    newg: newgame,
                });
            }


            taskcompleted = true;
            clearInterval(task25Id);
            timer25 = 1500;
            self.postMessage({
                event: 'setTick',
                minute: minute,
                second: second,
                moveYourAss: true,
                bringYourAssBack: false,
                timerLabel: 'Task',
            });
            minute = 0;
            second = 0;
            self.postMessage({
                event: 'metricUpdate',
                task: true,
                break: undefined,
            });
            // self.postMessage({event: 'mover25Completed'});
        }
    }, 1000);
}

function handleBreak() {


    if (taskcount !== 0 && taskcount % 5 !== 0 && taskcompleted) {

        // self.postMessage({event: 'tick5'});

        break5Id = setInterval(async function () {
            // this callback will execute every second until we call
            // the clearInterval method
            timer5--;
            console.log(timer5);

            second++;

            if (second > 59) {
                second = 0;
                minute++
            }

            self.postMessage({
                event: 'setTick',
                minute: minute,
                second: second,
                moveYourAss: false,
                bringYourAssBack: false,
                timerLabel: 'Break',
            });


            if (timer5 === 0) {
                // self.postMessage({event: 'break5'});
                clearInterval(break5Id);
                breakcount++

                lifecount = (100 * breakcount) / taskcount;

                self.postMessage({
                    event: 'scoreUpdate',
                    breaks: breakcount,
                    moves: taskcount,
                    lifes: lifecount,
                    newg: newgame,
                });

                timer5 = 300;

                self.postMessage({
                    event: 'setTick',
                    minute: minute,
                    second: second,
                    moveYourAss: false,
                    bringYourAssBack: true,
                    timerLabel: 'Break',
                });

                minute = 0;
                second = 0;

                self.postMessage({
                    event: 'metricUpdate',
                    task: undefined,
                    break: true,
                });
            }
        }, 1000);
    } else if (taskcount !== 0 && taskcount % 5 === 0 && taskcompleted) {


        // self.postMessage({event: 'tick30'});

        break30Id = setInterval(async function () {
            // this callback will execute every second until we call
            // the clearInterval method
            timer30--;
            console.log(timer30);

            second++;

            if (second > 59) {
                second = 0
                minute++
            }

            self.postMessage({
                event: 'setTick',
                minute: minute,
                second: second,
                moveYourAss: false,
                bringYourAssBack: false,
                timerLabel: 'Break',
            });


            if (timer30 === 0) {
                clearInterval(break30Id);
                breakcount++

                lifecount = (100 * breakcount) / taskcount;

                self.postMessage({
                    event: 'scoreUpdate',
                    breaks: breakcount,
                    moves: taskcount,
                    lifes: lifecount,
                    newg: newgame,
                });

                timer30 = 1800;

                // self.postMessage({event: 'break30'});

                self.postMessage({
                    event: 'setTick',
                    minute: minute,
                    second: second,
                    moveYourAss: false,
                    bringYourAssBack: true,
                    timerLabel: 'Break',
                });

                minute = 0;
                second = 0;

                self.postMessage({
                    event: 'metricUpdate',
                    task: undefined,
                    break: true,
                });
            }
        }, 1000);
    } else {
        // self.postMessage({event: 'break0'});

        self.postMessage({
            event: 'setTick',
            minute: undefined,
            second: undefined,
            moveYourAss: false,
            bringYourAssBack: true,
            timerLabel: 'Break',
        });


        minute = 0;
        second = 0;
    }

    if (newgame && taskcount === 1) {
        newgame = false;
    }
}