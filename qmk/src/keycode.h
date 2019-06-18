
#include <cstdio>
#include <stdint.h>

#include "keys.h"

using namespace std;

#define SEQ_IN_PROGRESS (1 << 0)
#define SEQ_COMPLETE (1 << 1)

#define SEQ_STATE_SET(bitfield) (seq_state |= bitfield)
#define SEQ_STATE_GET(bitfield) (seq_state & bitfield)
#define SEQ_STATE_CLEAR(bitfield) (seq_state &= ~bitfield)
#define SEQ_STATE_RESET() (seq_state = 0)

#define ID_KEY (0)
#define ID_LIST (1)
#define ID_MACRO (2)
#define ID_JUMP (3)

#define BITMASK (uint16_t(0x03fff))

#define SET_LIST(val) ((val & BITMASK) ^ (ID_LIST << 14))
#define SET_MACRO(val) ((val & BITMASK) ^ (ID_MACRO << 14))
#define SET_JUMP(val) ((val & BITMASK) ^ (ID_JUMP << 14))

#define GET_ID(val) (val >> 14)
#define GET_VAL(val) (val & BITMASK)

// #define IS_END(val) (val == 0)
// #define IS_KEY(val) (val>>14 == ID_KEY)
// #define IS_SIZE(val) (val>>14 == SIZE_ID)
// #define IS_MACRO(val) (val>>14 == ID_MACRO)
// #define IS_JUMP(val) (val>>14 == ID_JUMP)

const uint16_t seq_table[] = {SET_LIST(3), 'a', '1', 'n', '2', 'z', SET_JUMP(7), SET_LIST(2), 'c', '8', 'x', '9'};
uint16_t seq_pos = 1;
uint8_t seq_state = 0;
uint16_t seq_key = seq_table[0];
uint16_t seq_val = GET_VAL(seq_table[0]);
uint8_t seq_id = GET_ID(seq_table[0]);
uint16_t seq_macro = uint16_t(~0);

#define CODE_READ()                 \
	seq_key = seq_table[seq_pos++]; \
	seq_id = GET_ID(seq_key);       \
	seq_val = GET_VAL(seq_key);
#define CODE_JUMP(pos) (seq_pos = pos)
#define CODE_RESET() (seq_pos = 0)

#define seq_error(val, ...) ;

void seq_reset()
{
	seq_macro = 0;
	SEQ_STATE_RESET();
	CODE_RESET();
	CODE_READ();
}

void seq_do_result(void)
{
	printf("macro: %c (%d %d)\n", GET_VAL(seq_macro), GET_ID(seq_macro), GET_VAL(seq_macro));
}

bool seq_binary_search(uint16_t size, uint16_t keycode)
{
	uint16_t base_pos = seq_pos;
	int16_t l = 0;
	int16_t r = size - 1;
	printf("debug -> bin search l:%d r:%d\n", l, r);
	while (l <= r)
	{
		uint16_t m = (r + l) >> 1;
		CODE_JUMP((m << 1) + base_pos);
		// CODE_JUMP((r+1)|1);
		CODE_READ();
		printf("debug -> bin search got %c\n", seq_key);
		if (seq_key < keycode)
		{
			l = m + 1;
		}
		else if (seq_key > keycode)
		{
			r = m - 1;
		}
		else
		{
			// CODE_JUMP((m<<1)+1);
			// CODE_READ();
			printf("debug -> bin search code %d %d\n", seq_id, seq_val);
			// CODE_RESET();
			return true;
		}
	}
	return false;
}

bool seq_process(uint16_t keycode)
{
	printf("debug -> code %d %d\n", seq_id, seq_val);
	if (keycode == 0)
	{
		return true;
	}

	bool ok = true;

	if (seq_val == 0)
	{
		seq_error("zero at start", 0); // warning: seq_table is empty
		seq_reset();
		return false;
	}

	switch (seq_id)
	{
	default:
		seq_error("unexpected code id", 0); // error
		seq_reset();
		return false;
	case ID_KEY:
		ok = (seq_val == keycode);
		if (ok)
		{
			seq_macro = 0;
		}
		break;
	case ID_LIST:
		if (!seq_binary_search(seq_val, keycode))
		{
			ok = false;
			break;
		}
		seq_macro = 0;
		CODE_READ();
		switch (seq_id)
		{
		default:
			seq_error("wrong code id", 0); // error
			ok = false;
			break;
		case ID_KEY:
		case ID_MACRO:
			seq_macro = seq_key;
			SEQ_STATE_SET(SEQ_COMPLETE);
		case ID_JUMP:;
		}
	} // switch (seq_id)

	if (!ok)
	{
		if (SEQ_STATE_GET(SEQ_IN_PROGRESS))
		{
			if (seq_macro != 0)
			{
				seq_do_result();
			}
			else
			{
				ok = true;
			}
		}
		seq_reset();
		return ok;
	}
	SEQ_STATE_SET(SEQ_IN_PROGRESS);

	while (true)
	{

		if (seq_val == 0)
		{
			SEQ_STATE_SET(SEQ_COMPLETE);
		}

		if (SEQ_STATE_GET(SEQ_COMPLETE))
		{
			seq_do_result();
			seq_reset();
			return true;
		}

		switch (seq_id)
		{
		default:
			seq_error("unreachable", 0); // error
			seq_reset();
			return false;
		case ID_KEY:
		case ID_LIST:
			return true; // interrupt to wait next keycode
		case ID_JUMP:
			CODE_JUMP(seq_val);
			printf("debug -> jump to %d\n", seq_pos);
			break;
		case ID_MACRO:
			seq_macro = seq_val;
			break;
		} // switch (seq_id)
		CODE_READ();
	} // while (true)
}

bool process_record_user(uint16_t keycode)
{
	if (seq_process(keycode))
	{
		return false;
	}
	return true;
}