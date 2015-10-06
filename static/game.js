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

    function snake(x, y) {
        square(x, y, 'orange');
    }

    function food(x, y) {
        square(x, y, 'red');
    }

    function empty(x, y) {
        square(x, y, 'grey');
    }

    function square(x, y, color) {
        x++;
        y++;
        context.beginPath();
        // Top left
        context.moveTo(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP);
        // Top right
        context.lineTo(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP + SQUARE_WIDTH,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP);
        // Bottom right
        context.lineTo(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP + SQUARE_WIDTH,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP + SQUARE_HEIGHT);
        // Bottom left
        context.lineTo(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP + SQUARE_HEIGHT);
        // Back to the top left
        context.lineTo(
            x * SQUARE_WIDTH + x * SQUARE_HORIZONTAL_GAP,
            y * SQUARE_HEIGHT + y * SQUARE_VERTICAL_GAP);

        context.strokeStyle = 'black';
        context.fillStyle = color;
        context.fill();
    }

    function draw(board) {
        // board is a json structure.
        console.log(board);

        // Clear the canvas.
        context.save();
        context.setTransform(1, 0, 0, 1, 0, 0);
        context.clearRect(0, 0, canvas.width, canvas.height);
        context.restore();

        food(board.Food.X, board.Food.Y);
        snake(board.SnakeHead.X, board.SnakeHead.Y);
        snake(board.SnakeTail.X, board.SnakeTail.Y);
        for (var i = 0; i < board.Empty.length; i++) {
            var emptyPoint = board.Empty[i];
            //empty(emptyPoint.X, emptyPoint.Y);
            console.log(emptyPoint.X, emptyPoint.Y);
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
};
