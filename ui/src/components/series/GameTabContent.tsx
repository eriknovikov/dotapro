import * as React from "react"
import { Check, Clock, Copy, ExternalLink } from "lucide-react"
import type { SeriesMatchDetail, TeamInfo } from "@/api/api"
import { Button } from "@/components/ui"
import { PlayersTable } from "./index"
import { formatDuration, copyToClipboard } from "@/lib"

interface GameTabContentProps {
    match: SeriesMatchDetail
    gameNumber: number
    radiantTeam: TeamInfo
    direTeam: TeamInfo
}

export function GameTabContent({ match, gameNumber, radiantTeam, direTeam }: GameTabContentProps) {
    const [copied, setCopied] = React.useState(false)

    const handleCopyId = async () => {
        const success = await copyToClipboard(match.match_id.toString())
        if (success) {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        }
    }

    return (
        <div className="rounded-xl border border-white/10 bg-background-card focus:outline-none overflow-hidden" role="region" aria-label={`Game ${gameNumber} details`}>
            <div className="p-6">
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <div>
                        <h3 className="text-md" id={`match-${match.match_id}`}>Match ID: {match.match_id}</h3>
                        <p className="text-sm text-foreground-muted flex gap-1 items-center select-none">
                            <Clock size={15} aria-hidden="true" /> <span aria-label={`Duration: ${formatDuration(match.duration)}`}>{formatDuration(match.duration)}</span>
                        </p>
                    </div>
                    {/* Mobile: Icon-only buttons */}
                    <div className="flex gap-2 sm:hidden justify-center">
                        <Button variant="outline" size="sm" onClick={handleCopyId} className="p-2" aria-label="Copy ID">
                            {copied ? <Check className="w-4 h-4 text-green-600" /> : <Copy className="w-4 h-4" />}
                        </Button>

                        <Button variant="outline" size="sm" asChild className="p-2" aria-label="View on OpenDota">
                            <a
                                href={`https://opendota.com/matches/${match.match_id}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex gap-1 items-center"
                            >
                                OD
                                <ExternalLink className="w-3 h-3" aria-hidden="true" />
                            </a>
                        </Button>

                        <Button variant="outline" size="sm" asChild className="p-2" aria-label="View on Stratz">
                            <a
                                href={`https://stratz.com/matches/${match.match_id}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex gap-1 items-center"
                            >
                                Stratz
                                <ExternalLink className="w-3 h-3" aria-hidden="true" />
                            </a>
                        </Button>
                    </div>
                    {/* Desktop: Full text buttons */}
                    <div className="hidden sm:flex gap-2">
                        <Button variant="outline" size="sm" onClick={handleCopyId} className="gap-2 leading-relaxed">
                            {copied ? "Copied!" : "Copy ID"}
                            {!copied ? <Copy className="w-4 h-4" /> : <Check className="w-4 h-4 text-green-600" />}
                        </Button>

                        <Button variant="outline" size="sm" asChild>
                            <a
                                href={`https://opendota.com/matches/${match.match_id}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex gap-2 items-center leading-relaxed"
                                aria-label="View match on OpenDota"
                            >
                                <p>OpenDota</p>
                                <ExternalLink className="w-4 h-4" aria-hidden="true" />
                            </a>
                        </Button>
                        <Button variant="outline" size="sm" asChild>
                            <a
                                href={`https://stratz.com/matches/${match.match_id}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex gap-2 items-center leading-relaxed"
                                aria-label="View match on Stratz"
                            >
                                <p>Stratz</p>
                                <ExternalLink className="w-4 h-4" aria-hidden="true" />
                            </a>
                        </Button>
                    </div>
                </div>
            </div>
            <div className="px-6 pb-6">
                <PlayersTable match={match} radiantTeam={radiantTeam} direTeam={direTeam} />
            </div>
        </div>
    )
}
