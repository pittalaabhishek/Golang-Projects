package templates

const StoryTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Your Own Adventure - {{.Title}}</title>
    <style>
        body {
            font-family: 'Georgia', serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
            margin-bottom: 30px;
        }
        .story-content {
            margin-bottom: 30px;
        }
        .story-content p {
            margin-bottom: 20px;
            text-align: justify;
            font-size: 18px;
        }
        .options {
            margin-top: 40px;
        }
        .options h2 {
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .option {
            display: block;
            margin: 15px 0;
            padding: 15px 20px;
            background-color: #3498db;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            transition: background-color 0.3s ease;
            font-size: 16px;
        }
        .option:hover {
            background-color: #2980b9;
        }
        .end-message {
            text-align: center;
            font-style: italic;
            color: #7f8c8d;
            margin-top: 30px;
            padding: 20px;
            background-color: #ecf0f1;
            border-radius: 5px;
        }
        .restart {
            text-align: center;
            margin-top: 20px;
        }
        .restart a {
            color: #3498db;
            text-decoration: none;
            font-weight: bold;
        }
        .restart a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        
        <div class="story-content">
            {{range .Story}}
                <p>{{.}}</p>
            {{end}}
        </div>
        
        {{if .Options}}
            <div class="options">
                <h2>What do you choose?</h2>
                {{range .Options}}
                    <a href="/{{.Arc}}" class="option">{{.Text}}</a>
                {{end}}
            </div>
        {{else}}
            <div class="end-message">
                <h2>The End</h2>
                <p>You have reached the end of this story arc.</p>
            </div>
            <div class="restart">
                <a href="/intro">Start Over</a>
            </div>
        {{end}}
    </div>
</body>
</html>
`