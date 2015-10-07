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

    function snake(x, y) {
        square(x, y, 'orange');
    }

    function food(x, y) {
        square(x, y, 'red');
    }

    function empty(x, y) {
        square(x, y, 'grey');
    }

    function border(width, height) {
        context.beginPath();
        context.rect(0, 0, width, height);
        context.strokeStyle = 'black';
        context.stroke();
    }

    function draw(board) {
        // board is a json structure.
        console.log(board);

        // Clear the canvas.
        context.save();
        context.setTransform(1, 0, 0, 1, 0, 0);
        context.clearRect(0, 0, canvas.width, canvas.height);
        context.restore();

        border(
            board.Width * (SQUARE_WIDTH + SQUARE_HORIZONTAL_GAP) + SQUARE_WIDTH,
            board.Height * (SQUARE_HEIGHT + SQUARE_VERTICAL_GAP) + SQUARE_HEIGHT);

        food(board.Food.X, board.Food.Y);
        snake(board.SnakeHead.X, board.SnakeHead.Y);
        snake(board.SnakeTail.X, board.SnakeTail.Y);

        for (var x = 0; x < board.Width; x++) {
            for (var y = 0; y < board.Height; y++) {
                if (x == board.Food.X && y == board.Food.Y) {
                    continue;
                }
                var found = false;
                for (var i = 0; i < board.Empty.length; i++) {
                    var emptyPoint = board.Empty[i];
                    if (x == emptyPoint.X && y == emptyPoint.Y) {
                        found = true;
                        break;
                    }
                }
                if (!found) {
                    snake(x, y);
                }
            }
        }

    }

    function httpGet(url, callback) {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() {
            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                callback(xmlHttp.responseText);
            }
        };
        xmlHttp.open('GET', url, true);
        xmlHttp.send(null);
    }

    var left = document.getElementById('left');
    var right = document.getElementById('right');
    var forward = document.getElementById('forward');
    var reset = document.getElementById('reset');

    left.onclick = function() {
        httpGet('snake/move/left', function(text) {
            draw(JSON.parse(text));
        });
    };

    right.onclick = function() {
        httpGet('snake/move/right', function(text) {
            draw(JSON.parse(text));
        });
    };

    forward.onclick = function() {
        httpGet('snake/move/forward', function(text) {
            draw(JSON.parse(text));
        });
    };

    reset.onclick = function() {
        httpGet('snake/reset', function(text) {
            draw(JSON.parse(text));
        });
    };

    reset.onclick();
};
