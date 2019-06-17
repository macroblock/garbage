
#include <cstdio>

#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>

#include "bitmapfont.h"

font_t* newFont(SDL_Renderer* renderer, const char* filename, int first, int last, int perLine, int runeW, int runeH, int dX, int dY) {
    auto tex = IMG_LoadTexture(renderer, filename);
	if (tex == NULL) {
		fprintf(stderr, "load texture error\n");
		fprintf(stderr, "%s\n", SDL_GetError());

		return NULL;
	}
	SDL_SetTextureBlendMode(tex, SDL_BLENDMODE_BLEND);

	font_t* font = new font_t;
	font->renderer = renderer;
	font->tex = tex;
	font->firstRune = first;
	font->lastRune = last;
	font->runesPerLine = perLine;
	font->runeW = runeW;
	font->runeH = runeH;
	font->dX = dX;
	font->dY = dY;
	font->scale = 1;
	return font;
}

void freeFont(font_t** font) {
	if (font == NULL) {
		return;
	}
	if ((*font)->tex != NULL) {
		SDL_DestroyTexture( (*font)->tex );
		(*font)->tex = NULL;
	}
	delete (*font);
	*font = NULL;
}

void setFontScale(font_t* font, int scale) {
	if (font == NULL) {
		return;
	}
	font->scale = scale;
}

void setFontColor(font_t* font, color_t color) {
	if (font == NULL) {
		return;
	}
	SDL_SetTextureColorMod(font->tex, color.r, color.g, color.b);
	SDL_SetTextureAlphaMod(font->tex, color.a);
}

void printRune(font_t* font, int x, int y, int rune) {
	if (font == NULL) {
		return;
	}
	if (font->firstRune > rune || font->lastRune < rune) {
		rune = font->lastRune + 1; // rune substitude
	}
	auto r = rune - font->firstRune;
	auto rx = r % font->runesPerLine;
	auto ry = r / font->runesPerLine;
	auto rw = font->runeW;
	auto rh = font->runeH;
	auto src = SDL_Rect{rx*rw, ry*rh, rw, rh};
	auto dst = SDL_Rect{x, y, rw*font->scale, rh*font->scale};
	SDL_RenderCopy(font->renderer, font->tex, &src, &dst);
}

void printText(font_t* font, int x0, int y0, const char* text) {
	if (font == NULL) {
		return;
	}
	auto x = x0, y = y0;
	for (int i = 0; text[i] != '\0'; i++) {
		auto r = text[i];
		// printf("%d %c\n", i, char(r));
		switch (r) {
			case '\n':
			case '\r':
				x = x0;
				y += (font->runeH + font->dY)*font->scale;
				continue;
		}
		printRune(font, x, y, r);
		x += (font->runeW + font->dX)*font->scale;
	}
}