#include <cstdio>
#include <stdint.h>

#include <string>

// #include <ncurses.h>
// #define printf(...) printw(__VA_ARGS__); refresh()

#include "keycode.h"

#include "bitmapfont.h"

// int main()
// {
// 	initscr();
// 	// raw();
// 	// def_prog_mode();
// 	scrollok(stdscr,TRUE);

//     addstr("type a key\n");
//     while (1) {
// 		refresh();
//         uint16_t key = getch();
// 		if (key == 65535) {
// 			printf("error: getch()\n");
// 			break;
// 		}
//         printf("\n> keycode: %c, %d\n", key, key);
//         if (key == 0x1b) {
//             break;
//         }
//         if (process_record_user(key)) {
//             printf("unhandled keycode: %c, %d\n", key, key);
//         }
//     }
// 	endwin();
// }

// #include <stdio.h> /* printf and fprintf */

#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>

#define WIDTH 800
#define HEIGHT 600

color_t colorDefault = {200,200,200,255};
color_t colorError = {255,0,0,255};
color_t colorDown = {50,200,50,255};
color_t colorUp = {100,100,255,255};
color_t colorPress = {50,200,200,255};
color_t colorRepeat = {200,200,50,255};
color_t colorMod = {200,0,200,255};

void keySeqPrint(font_t* font, int x0, int y0, const char* text) {
	if (font == NULL) {
		return;
	}
	int w, h;
	SDL_GetRendererOutputSize(font->renderer, &w, &h);
	setFontColor(font, colorDefault);
	auto readParameter = false;
	auto x = x0, y = y0;
	for (int i = 0; text[i] != '\0'; i++) {
		auto r = text[i];
		// printf("%d %c\n", i, char(r));
		if (readParameter) {
			readParameter = false;
			switch (r) {
			default:
				setFontColor(font, colorError);
				printRune(font, x, y, -1);
				setFontColor(font, colorDefault);
				break;
			case 'd':
				setFontColor(font, colorDown);
				break;
			case 'u':
				setFontColor(font, colorUp);
				break;
			case 'p':
				setFontColor(font, colorPress);
				break;
			case 'r':
				setFontColor(font, colorRepeat);
				break;
			case 'm':
				setFontColor(font, colorMod);
				break;
			case '0':
				setFontColor(font, colorDefault);
				break;
			}
			continue;
		} else {
			switch (r) {
			case '\n':
			case '\r':
				x = x0;
				y += (font->runeH + font->dY)*font->scale;
				continue;
			case 0x1b:
				readParameter = true;
				continue;
			}
			printRune(font, x, y, r);
			x += (font->runeW + font->dX)*font->scale;
			if (x > w) {
				x = x0;
				y += (font->runeH + font->dY)*font->scale;
			}
			if (y > h) {
			}
		}
	}
}

void keySeqAppend(std::string* str, bool pressed, int repeat, SDL_Keysym keysym) {
	auto key = SDL_GetKeyFromScancode(keysym.scancode);
	auto keyName = SDL_GetKeyName(key);
	auto prefix = pressed ? "\x1b" "d" : "\x1b" "u";
	prefix = repeat>0 ? "\x1b" "r" : prefix;
	str->append(prefix);
	if (32 <= key && key <= 128) {
		str->append(1, char(key));
		return;
	}
	str->append("<");
	str->append(keyName);
	str->append(">");
}

#ifdef _WIN32
int WinMain(int, char**)
#else
int main()
#endif
{
	SDL_Window* window = NULL;
	SDL_Renderer* renderer = NULL;

	if (SDL_Init(SDL_INIT_VIDEO) != 0)
	{
		fprintf(stderr, "SDL failed to initialise: %s\n", SDL_GetError());
		return 1;
	}

	window = SDL_CreateWindow("qmk-test",			   /* Title of the SDL window */
							  SDL_WINDOWPOS_UNDEFINED, /* Position x of the window */
							  SDL_WINDOWPOS_UNDEFINED, /* Position y of the window */
							  WIDTH,				   /* Width of the window in pixels */
							  HEIGHT,				   /* Height of the window in pixels */
							  0);					   /* Additional flag(s) */
	if (window == NULL)
	{
		fprintf(stderr, "SDL window failed to initialise: %s\n", SDL_GetError());
		return 1;
	}

    renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);

	if (renderer == NULL)
	{
		fprintf(stderr, "SDL renderer failed to initialise: %s\n", SDL_GetError());
		return 1;
	}

	IMG_Init(IMG_INIT_PNG);

	font_t* font = newFont(renderer, "font-5x9-(16x7).png",
							32, 127, 16, // first, last, per line
							5, 9,  // rune size
							1, 2); // dx, dy
	setFontScale(font, 2);
	// fprintf(stderr, "temp %d\n", font->lastRune);
	// fprintf(stderr, "temp %p\n", (void*)font);

	SDL_Event event;
	bool quit = false;
	bool refresh = true;
	std::string text;
	while (!quit)
	{
		while (SDL_PollEvent(&event))
		{
			refresh = true;
			/* an event was found */
			switch (event.type)
			{
			case SDL_QUIT:
				quit = 1;
				break;
			case SDL_KEYDOWN:
			case SDL_KEYUP:
				auto pressed = (event.key.state == SDL_PRESSED);
				switch (event.key.keysym.sym)
				{
				default:
					keySeqAppend(&text, pressed, event.key.repeat, event.key.keysym);
					break;
				case SDLK_ESCAPE:
					quit = true;
					break;
				case SDLK_RETURN:
					if (pressed) {
						text.append(1, 13);
					}
					break;
				}
				break;
			}
		}
		if (refresh) {
			SDL_RenderClear(renderer);
			// SDL_RenderCopy(renderer, font->tex, NULL, NULL);
			keySeqPrint(font, 1, 1, text.c_str());
        	SDL_RenderPresent(renderer);
		}
		SDL_Delay(10);
	}

	freeFont(&font);

	SDL_DestroyWindow(window);
	SDL_Quit();
	return 0;
}
