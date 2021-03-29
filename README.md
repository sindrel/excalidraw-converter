# Excalidraw Converter
**A command line tool for porting Excalidraw diagrams to Gliffy.**

[Excalidraw](https://excalidraw.com/) is great for sketching diagrams as part of a design process, but chances are that you have to redo those sketches for documentation. This tool is made to bridge those tasks.

Excalidraw Converter ports Excalidraw diagrams to a Gliffy compatible format, which can be imported directly into services like [Gliffy](https://www.gliffy.com/) or [Gliffy for Confluence](https://marketplace.atlassian.com/apps/254/gliffy-diagrams-for-confluence).

![Excalidraw vs. Gliffy comparison](exconv-comparison.png "Comparison")

## Getting started

### Installation

Download a compatible binary from the [Releases](https://github.com/sindrel/excalidraw-converter/releases) page.

If you're a Linux or MacOS user, move it to your local bin folder to make it available in your environment (optional).

### How to convert diagrams
First save your Excalidraw diagram to a file.

Then, to do a conversion, simply execute the binary by specifying the `gliffy` command, the path to your Excalidraw save file, and the path to where you want your converted file to be saved. 

<details>
  <summary>MacOS example</summary>
  
  ```
  $ ./exconv gliffy ~/Downloads/my-diagram.excalidraw /tmp/my-ported-diagram.gliffy
  Parsing input file: ~/Downloads/my-diagram.excalidraw
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.line
  Converted diagram saved to file: /tmp/my-ported-diagram.gliffy
  ```
</details>

<details>
  <summary>Linux example</summary>
  
  ```
  $ ./exconv gliffy ~/Downloads/my-diagram.excalidraw /tmp/my-ported-diagram.gliffy
  Parsing input file: ~/Downloads/my-diagram.excalidraw
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.line
  Converted diagram saved to file: /tmp/my-ported-diagram.gliffy
  ```
</details>

<details>
  <summary>Windows example</summary>
  
  ```
  C:\> exconv.exe gliffy C:\Downloads\my-diagram.excalidraw C:\tmp\my-ported-diagram.gliffy
  Parsing input file: C:\Downloads\my-diagram.excalidraw
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.rectangle
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.text
  Adding object: com.gliffy.shape.basic.basic_v1.default.line
  Converted diagram saved to file: C:\tmp\my-ported-diagram.gliffy
  ```
</details>

![Animation demonstrating use](exconv.gif "Animation")

After converting your diagram(s), import them into Gliffy using the standard Import dialog.

## Commands

| Command | Use | Arguments |
|------|-------------|:----:|
| gliffy | Converts diagram to Gliffy format | \<input file path\> \<output file path\> |
| help | Shows help | N/A |
| version | Outputs the application version | N/A |

## Features

All fixed shapes and most styling and text options are supported.

### Shapes
* Rectangle
* Diamond
* Ellipse
* Arrow
* Line

### Text

* Font family (Normal and Code)
* Font size
* Font color
* Horizontal alignment

### Styling

* Canvas background color 
* Fill color
* Fill style (hachure and cross-hatch translate to gradients) 
* Stroke color
* Stroke width
* Opacity

Free hand drawings and library graphics are currently not supported.

## Contributing

See something you'd like to improve? Feel free to add a pull request. If it's a major change, it's probably best to describe it in an [issue](https://github.com/sindrel/excalidraw-converter/issues/new) first.