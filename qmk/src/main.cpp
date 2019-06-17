#include <cstdio>
#include <stdint.h>

// #include <ncurses.h>
// #define printf(...) printw(__VA_ARGS__); refresh()

#include "keycode.h"

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

#ifdef _WIN32
#include <SDL/SDL.h> /* Windows-specific SDL2 library */
#else
#include <SDL2/SDL.h> /* macOS- and GNU/Linux-specific */
#endif

/* Sets constants */
#define WIDTH 800
#define HEIGHT 600
#define DELAY 3000

int main()
{
	/* Initialises data */
	SDL_Window *window = NULL;

	/*
  * Initialises the SDL video subsystem (as well as the events subsystem).
  * Returns 0 on success or a negative error code on failure using SDL_GetError().
  */
	if (SDL_Init(SDL_INIT_VIDEO) != 0)
	{
		fprintf(stderr, "SDL failed to initialise: %s\n", SDL_GetError());
		return 1;
	}

	/* Creates a SDL window */
	window = SDL_CreateWindow("SDL Example",		   /* Title of the SDL window */
							  SDL_WINDOWPOS_UNDEFINED, /* Position x of the window */
							  SDL_WINDOWPOS_UNDEFINED, /* Position y of the window */
							  WIDTH,				   /* Width of the window in pixels */
							  HEIGHT,				   /* Height of the window in pixels */
							  0);					   /* Additional flag(s) */

	/* Checks if window has been created; if not, exits program */
	if (window == NULL)
	{
		fprintf(stderr, "SDL window failed to initialise: %s\n", SDL_GetError());
		return 1;
	}

	SDL_Event event;
	int gameover = 0;

	/* message pump */
	while (!gameover)
	{
		/* look for an event */
		if (SDL_PollEvent(&event))
		{
			/* an event was found */
			switch (event.type)
			{
			/* close button clicked */
			case SDL_QUIT:
				gameover = 1;
				break;

			/* handle the keyboard */
			case SDL_KEYDOWN:
				switch (event.key.keysym.sym)
				{
				case SDLK_ESCAPE:
				case SDLK_q:
					gameover = 1;
					break;
				}
				break;
			}
		}
	}

	//   /* Pauses all SDL subsystems for a variable amount of milliseconds */
	//   SDL_Delay(DELAY);

	/* Frees memory */
	SDL_DestroyWindow(window);

	/* Shuts down all SDL subsystems */
	SDL_Quit();

	return 0;
}