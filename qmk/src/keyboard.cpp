
#include <SDL2/SDL.h>

#include "keyboard.h"
#include "qmk.h"

#define SC(key) (SDL_SCANCODE_##key)
#define SCXXX (SDL_SCANCODE_UNKNOWN)

void qmk_init_keyboard(void) {
    init_scancode_table( MATRIX_COLS, MATRIX_ROWS,
SC(ESCAPE), SC(F1),  SC(F2),  SC(F3),   SC(F4),  SC(F5),  SC(F6),  SC(F7),   SC(F8),   SC(F9),    SC(F10),      SC(F11),        SC(F12),         SCXXX,        SCXXX,
SC(GRAVE),  SC(1),   SC(2),   SC(3),    SC(4),   SC(5),   SC(6),   SC(7),    SC(8),    SC(9),     SC(0),        SC(MINUS),      SC(EQUALS),      SC(BACKSPACE),SCXXX,
SC(TAB),    SC(Q),   SC(W),   SC(E),    SC(R),   SC(T),   SC(Y),   SC(U),    SC(I),    SC(O),     SC(P),        SC(LEFTBRACKET),SC(RIGHTBRACKET),SC(BACKSLASH),SCXXX,
SC(LCTRL),  SC(A),   SC(S),   SC(D),    SC(F),   SC(G),   SC(H),   SC(J),    SC(K),    SC(L),     SC(SEMICOLON),SC(APOSTROPHE), SC(RETURN),      SCXXX,        SCXXX,
SC(LSHIFT), SC(Z),   SC(X),   SC(C),    SC(V),   SC(B),   SC(N),   SC(M),    SC(COMMA),SC(PERIOD),SC(SLASH),    SC(RSHIFT),     SCXXX,           SCXXX,        SCXXX,
SC(LCTRL),  SC(LGUI),SC(LALT),SC(SPACE),SC(RALT),SC(RGUI),SC(MENU),SC(RCTRL),SCXXX,    SCXXX,     SCXXX,        SCXXX,          SCXXX,           SCXXX,        SCXXX
    );
}
