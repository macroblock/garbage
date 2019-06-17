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
	setFontScale(font, 3);
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
				switch (event.key.keysym.sym)
				{
				default:
					setFontColor(font, color_t{255,255,0, 127});
					text.append(1, char(event.key.keysym.sym));
					printText(font, 1, 1, text.c_str());
					break;
				case SDLK_ESCAPE:
				// case SDLK_q:
					quit = true;
					break;
				}
				break;
			case SDL_KEYUP:
				setFontColor(font, color_t{0,255,255, 64});
				text.append(1, char(event.key.keysym.sym));
				printText(font, 1, 1, text.c_str());
				break;
			}
		}
		if (refresh) {
			// SDL_RenderCopy(renderer, font->tex, NULL, NULL);
        	SDL_RenderPresent(renderer);
		}
		SDL_Delay(10);
	}

	freeFont(&font);

	SDL_DestroyWindow(window);
	SDL_Quit();
	return 0;
}
