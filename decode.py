import peachpy.x86_64

dst = Argument(ptr())
dst_len = Argument(int64_t)
dst_cap = Argument(int64_t)
src = Argument(ptr())
src_len = Argument(int64_t)
src_cap = Argument(int64_t)

with Function("Decode4", (dst, dst_len, dst_cap, src, src_len, src_cap), target=uarch.default + isa.sse4_1) as function:
    reg_dst = GeneralPurposeRegister64()
    reg_src = GeneralPurposeRegister64()
    reg_masks = GeneralPurposeRegister64()

    LOAD.ARGUMENT(reg_dst, dst)
    LOAD.ARGUMENT(reg_src, src)
    LOAD.ARGUMENT(reg_masks, dst_len)

    xmm = XMMRegister()
    MOVDQU(xmm, [reg_src+1])

    mask = GeneralPurposeRegister64()
    MOVZX(mask, byte[reg_src])
    SHL(mask, 4)
    ADD(mask, reg_masks)

    pshubfmask = XMMRegister()
    MOVDQA(pshubfmask, [mask])
    PSHUFB(xmm, pshubfmask)

    MOVDQU([reg_dst], xmm)

    RETURN()
