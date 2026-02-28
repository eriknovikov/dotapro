import { getItemById, getNeutralById, getItemByName } from "@/lib/dotautils"
import { getItemImageUrl } from "@/lib"

export function ItemTooltipContent({ itemId }: { itemId: number }) {
    const item = getItemById(itemId) || getNeutralById(itemId)
    if (!item || item.id === 0) return null

    return (
        <div className="w-72 bg-[#1c1d21] text-[#a1aab3] text-sm border border-[#30303f] shadow-xl rounded-md overflow-hidden font-sans focus-visible:outline-none">
            {/* Header */}
            <div className="flex items-center gap-3 p-3 bg-linear-to-r from-[#282836] to-[#1c1d21] border-b border-[#30303f]">
                <img src={getItemImageUrl(item.id)!} alt={item.displayName} className="w-14 h-10 rounded shadow-md" />
                <div>
                    <div className="text-white font-bold text-base">{item.displayName}</div>
                    {item.tier ? (
                        <div className="text-[#e5c158] text-xs font-semibold">Tier {item.tier} Neutral</div>
                    ) : (
                        <div className="flex items-center gap-1 text-[#e5c158] font-semibold">
                            <img
                                src="https://cdn.cloudflare.steamstatic.com/apps/dota2/images/tooltips/gold.png"
                                alt="Gold"
                                className="w-4 h-4"
                            />
                            {item.cost}
                        </div>
                    )}
                </div>
            </div>

            <div className="p-3 flex flex-col gap-3">
                {/* Attributes */}
                {item.attrib.length > 0 && (
                    <div className="flex flex-col gap-0.5">
                        {item.attrib.map((attr, idx) => {
                            if (!attr.display) return null
                            const parts = attr.display.split("{value}")
                            return (
                                <div key={idx}>
                                    {parts[0]}
                                    <span className="text-white font-semibold">{attr.value}</span>
                                    {parts[1]}
                                </div>
                            )
                        })}
                    </div>
                )}

                {/* Abilities */}
                {item.abilities.map((ability, idx) => (
                    <div key={idx} className="bg-[#141518] rounded border border-[#23242b] overflow-hidden">
                        <div className="bg-[#23242b] px-2 py-1 flex justify-between items-center text-white font-semibold">
                            <span>
                                {ability.type === "active" ? "Active: " : "Passive: "}
                                {ability.title}
                            </span>
                            <div className="flex gap-2 text-xs font-normal">
                                {item.cd && <span className="flex items-center gap-1 text-[#95c07a]">⏱ {item.cd}</span>}
                                {item.mc && (
                                    <span className="flex items-center gap-1 text-[#69a1e2]">💧 {item.mc}</span>
                                )}
                            </div>
                        </div>
                        <div className="p-2 text-[#8c9aa6] whitespace-pre-wrap text-xs">{ability.description}</div>
                    </div>
                ))}

                {/* Components */}
                {item.components && item.components.length > 0 && (
                    <div className="pt-2 border-t border-[#30303f]">
                        <div className="text-xs mb-1 text-white">Components:</div>
                        <div className="flex flex-wrap gap-1">
                            {item.components.map((compName, idx) => {
                                const compItem = getItemByName(compName)
                                if (!compItem) return null
                                return (
                                    <div key={idx} className="relative group">
                                        <img
                                            src={getItemImageUrl(compItem.id)!}
                                            alt={compItem.displayName}
                                            className="w-10 h-7 rounded"
                                        />
                                        <div className="absolute bottom-0 right-0 bg-black/80 text-[#e5c158] text-[10px] px-1 rounded-tl">
                                            {compItem.cost}
                                        </div>
                                    </div>
                                )
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
