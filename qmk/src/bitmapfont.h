#ifndef __BITMAPFONT_H__
#define __BITMAPFONT_H__

#include <stdint.h>

#include <SDL2/SDL.h>

typedef struct {
	int firstRune;
	int lastRune;
	int runesPerLine;
	int runeW, runeH;
	int dX, dY;
	int scale;
	SDL_Renderer* renderer;
	SDL_Texture* tex;
} font_t;

typedef struct {
    uint8_t r, g, b, a;
} color_t;

font_t* newFont(SDL_Renderer* renderer, const char* filename, int first, int last, int perLine, int runeW, int runeH, int dX, int dY);

void freeFont(font_t** font);

void printText(font_t* font, int x0, int y0, const char* text);
void printRune(font_t* font, int x, int y, int rune);

void setFontScale(font_t* font, int scale);
void setFontColor(font_t* font, color_t color);


#endif // __BITMAPFONT_H__