# Conway's Game of Life

An interactive implementation of [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) built in
Go with SDL2.

https://github.com/user-attachments/assets/d1ff7345-0b97-4a6e-8a93-6cc05ca48947

## Quick Start

Download a pre-built binary from [Releases](../../releases).

**macOS (Apple Silicon):**

```bash
tar xzf gameoflife-darwin-arm64.tar.gz
xattr -d com.apple.quarantine gameoflife-darwin-arm64
./gameoflife-darwin-arm64
```

**Linux (amd64):**

```bash
# Install SDL2 runtime libraries if not already present
sudo apt-get install libsdl2-2.0-0 libsdl2-ttf-2.0-0   # Debian/Ubuntu
sudo dnf install SDL2 SDL2_ttf                            # Fedora

tar xzf gameoflife-linux-amd64.tar.gz
./gameoflife-linux-amd64
```

## Features

- Infinite sparse grid (no fixed boundaries)
- 15 built-in classic patterns (Glider Gun, Pulsar, Acorn, R-pentomino, and more)
- Heatmap mode: color cells by age (green -> yellow -> red -> purple)
- Fade trails: dead cells leave a fading ghost trail
- Adjustable simulation speed (1-60 generations/sec)
- Smooth zoom and pan (mouse + keyboard)
- Fullscreen mode
- Draw/erase cells with mouse
- Load custom patterns from RLE and Golly Macrocell (.mc) files

## Controls

| Key                | Action                                  |
|--------------------|-----------------------------------------|
| `P` / `Space`      | Play / Pause                            |
| `S`                | Step one generation (while paused)      |
| `G`                | Toggle grid overlay                     |
| `H`                | Toggle heatmap mode (color by cell age) |
| `T`                | Toggle fade trails                      |
| `F`                | Toggle fullscreen                       |
| `C`                | Clear the grid                          |
| `R`                | Reset to current pattern                |
| `N`                | Next pattern from built-in catalog      |
| `[` / `]`          | Decrease / increase simulation speed    |
| `+` / `-` / Scroll | Zoom in / out                           |
| Arrow keys         | Pan camera                              |
| Left click         | Toggle cell alive/dead                  |
| Left drag          | Draw cells                              |
| Right drag         | Pan camera                              |
| `ESC`              | Exit                                    |

## Built-in Patterns

Press `N` to cycle through these while running:

| Pattern               | Type       | Description                                              |
|-----------------------|------------|----------------------------------------------------------|
| R-pentomino           | Methuselah | 5 cells, stabilizes after 1103 generations               |
| Glider                | Spaceship  | The smallest known spaceship                             |
| Gosper Glider Gun     | Gun        | First known finite pattern with infinite growth          |
| Lightweight Spaceship | Spaceship  | Travels horizontally                                     |
| Heavyweight Spaceship | Spaceship  | Larger horizontal spaceship                              |
| Pulsar                | Oscillator | Period 3, the most common naturally occurring oscillator |
| Pentadecathlon        | Oscillator | Period 15                                                |
| Beacon                | Oscillator | Period 2                                                 |
| Blinker               | Oscillator | Period 2, simplest oscillator                            |
| Toad                  | Oscillator | Period 2                                                 |
| Block                 | Still life | 2x2 stable pattern                                       |
| Beehive               | Still life | 6-cell stable pattern                                    |
| Acorn                 | Methuselah | 7 cells, runs for 5206 generations                       |
| Diehard               | Methuselah | 7 cells, dies completely after 130 generations           |
| Infinite Growth       | Infinite   | 1-row pattern that grows forever                         |

## Pattern Files

The `patterns/` directory includes RLE files for individual components and a Macrocell file for a full circuit:

| File                       | Component            | Description                                          |
|----------------------------|----------------------|------------------------------------------------------|
| `glider.rle`               | Glider               | Signal carrier in digital logic (c/4 diagonal)       |
| `gosper-glider-gun.rle`    | Gosper Glider Gun    | Clock signal / constant "1" source (period-30)       |
| `eater1.rle`               | Eater 1              | Signal sink — absorbs gliders without damage         |
| `buckaroo.rle`             | Buckaroo             | Period-30 glider reflector (90-degree wire bend)     |
| `snark.rle`                | Snark                | Smallest known stable reflector (43-gen recovery)    |
| `inline-inverter.rle`      | NOT Gate             | Inverts a period-30 glider stream                    |
| `glider-duplicator.rle`    | Fanout / Duplicator  | Splits one signal into two output streams            |
| `r-pentomino.rle`          | R-pentomino          | Classic methuselah (1103 generations to stabilize)   |

## Custom Patterns

Load [RLE](https://conwaylife.com/wiki/Run_Length_Encoded) or [Macrocell](https://conwaylife.com/wiki/Macrocell) pattern files:

```bash
./gameoflife path/to/pattern.rle
./gameoflife path/to/pattern.mc
```

Thousands of patterns are available at [LifeWiki](https://conwaylife.com/wiki/Main_Page).

## Building from Source

### Prerequisites

Install SDL2 and SDL2_ttf development libraries:

**macOS (Homebrew):**

```bash
brew install sdl2 sdl2_ttf
```

**Ubuntu/Debian:**

```bash
sudo apt-get install libsdl2-dev libsdl2-ttf-dev
```

**Fedora:**

```bash
sudo dnf install SDL2-devel SDL2_ttf-devel
```

### Build

```bash
# Build and run (default pattern: R-pentomino)
make run

# Run with a specific pattern file
make build
bin/gameoflife patterns/gosper-glider-gun.rle

# Build with race detector (development)
make run_dev
```

## License

MIT
