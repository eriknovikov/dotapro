import { getHeroById } from "@/lib/dotautils"
import { getHeroImageUrl } from "@/lib"

export function HeroTooltipContent({ heroId }: { heroId: number }) {
    const hero = getHeroById(heroId)
    if (!hero) return null

    return (
        <div className="w-80 bg-[#1c1d21] text-[#a1aab3] text-sm border border-[#30303f] shadow-xl rounded-md overflow-hidden font-sans">
            {/* Header */}
            <div className="flex items-center gap-3 p-3 bg-linear-to-r from-[#282836] to-[#1c1d21] border-b border-[#30303f]">
                <img src={getHeroImageUrl(hero.id)} alt="" className="w-16 h-10.5 rounded shadow-md" />
                <div>
                    <div className="text-white font-bold text-base">{hero.displayName}</div>
                </div>
            </div>
            {/* Lore */}
            <div className="p-3">
                <p className="text-xs leading-relaxed">{hero.lore}</p>
            </div>
        </div>
    )
}
