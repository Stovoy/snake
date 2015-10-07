document.onreadystatechange = function () {
    if (document.readyState != 'complete') {
        return;
    }
    
    var canvas = document.getElementById('canvas');
    var context = canvas.getContext('2d');

    canvas.width = 800;
    canvas.height = 800;

    var SQUARE_WIDTH = 4;
    var SQUARE_HEIGHT = 4;
    var SQUARE_HORIZONTAL_GAP = 2;
    var SQUARE_VERTICAL_GAP = 2;
    var PADDING = 2;

    function square(x, y, color) {
        context.fillStyle = color;
        context.fillRect(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP + PADDING,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP + PADDING,
            SQUARE_WIDTH, SQUARE_HEIGHT);
    }

    /**
     * Connect two adjacent squares.
     * @param x1
     * @param y1
     * @param x2
     * @param y2
     * @param color
     */
    function connect(x1, y1, x2, y2, color) {
        context.fillStyle = color;
        var lower;
        if (Math.abs(x1 - x2) === 1) {
            if (y1 != y2) {
                return;
            }
            lower = Math.min(x1, x2);
            context.fillRect(
                lower * SQUARE_WIDTH + lower * SQUARE_HORIZONTAL_GAP + PADDING + SQUARE_WIDTH,
                y1 * SQUARE_HEIGHT + y1 * SQUARE_VERTICAL_GAP + PADDING,
                SQUARE_HORIZONTAL_GAP, SQUARE_HEIGHT);
        } else if (Math.abs(y1 - y2) === 1) {
            if (x1 != x2) {
                return;
            }
            lower = Math.min(y1, y2);
            context.fillRect(
                x1 * SQUARE_WIDTH + x1 * SQUARE_HORIZONTAL_GAP + PADDING,
                lower * SQUARE_HEIGHT + lower * SQUARE_VERTICAL_GAP + PADDING + SQUARE_HEIGHT,
                SQUARE_WIDTH, SQUARE_VERTICAL_GAP);
        }
    }

    /**
     * List of snake positions, in order. Must be connected.
     * @param positions
     */
    function snake(positions) {
        var previous = null;
        for (var i = 0; i < positions.length; i++) {
            var position = positions[i];
            square(position.x, position.y, 'orange');
            if (previous !== null) {
                connect(previous.x, previous.y, position.x, position.y, 'orange');
            }
            previous = position;
        }
    }

    function food(x, y) {
        square(x, y, 'red');
    }

    function border(width, height) {
        context.beginPath();
        context.rect(0, 0, width, height);
        context.strokeStyle = 'black';
        context.stroke();
    }

    function draw(board) {
        // board is a json structure.

        // Clear the canvas.
        context.save();
        context.setTransform(1, 0, 0, 1, 0, 0);
        context.clearRect(0, 0, canvas.width, canvas.height);
        context.restore();

        border(
            board.Width * (SQUARE_WIDTH + SQUARE_HORIZONTAL_GAP) + SQUARE_WIDTH,
            board.Height * (SQUARE_HEIGHT + SQUARE_VERTICAL_GAP) + SQUARE_HEIGHT);

        food(board.Food.X, board.Food.Y);
        var snakeList = [{x: board.SnakeHead.X, y: board.SnakeHead.Y}];

        for (var i = board.SnakeBody.length - 1; i >= 0; i--) {
            var body = board.SnakeBody[i];
            snakeList.push({x: body.X, y: body.Y})
        }

        snake(snakeList);
    }

    function httpGet(url, callback) {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() {
            if (xmlHttp.readyState === 4 && xmlHttp.status === 200) {
                callback(xmlHttp.responseText);
            }
        };
        xmlHttp.open('GET', url, true);
        xmlHttp.send(null);
    }

    function drawCallback(text) {
        draw(JSON.parse(text));
    }

    function left() {
        httpGet('snake/move/left', drawCallback);
    }

    function right() {
        httpGet('snake/move/right', drawCallback);
    }

    function forward() {
        httpGet('snake/move/forward', drawCallback);
    }

    function reset() {
        httpGet('snake/reset', drawCallback);
    }

    function rewind() {
        httpGet('snake/rewind', drawCallback);
    }

    var leftButton = document.getElementById('left');
    var rightButton = document.getElementById('right');
    var forwardButton = document.getElementById('forward');
    var resetButton = document.getElementById('reset');
    var rewindButton = document.getElementById('rewind');

    leftButton.onclick = left;
    rightButton.onclick = right;
    forwardButton.onclick = forward;
    resetButton.onclick = reset;
    rewindButton.onclick = rewind;

    reset();

    document.onkeypress = function (e) {
        e = e || window.event;
        if (e.keyCode === 119) {
            // w
            forward();
        } else if (e.keyCode == 97) {
            // a
            left();
        } else if (e.keyCode == 100) {
            // d
            right();
        } else if (e.keyCode == 114) {
            // r
            rewind();
        }
    };
};
