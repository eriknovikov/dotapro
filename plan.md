# Implementation Plan: Enhanced Team Header in PlayersTable

## Overview

Enhance the `PlayersTable.tsx` component to display a more informative team header with:

- Full team name
- Team score
- "WON" badge for the winning team (with Trophy icon)

## Confirmed Decisions

### 1. Data Flow: Prop Drilling (Option A)

Pass team info from parent component through the component chain:

- `SeriesDetailPage` (`$id.tsx`) → `GameTabContent` → `PlayersTable`

**Rationale:** No API changes required, manageable prop chain.

### 2. Team Icons: Separate Component (Option C)

Create a reusable `TeamIcon` component that accepts SVG paths for Radiant and Dire.

**Rationale:** User will provide suitable SVGs for both teams. Reusable and clean.

### 3. "WON" Badge: Custom Component with Trophy Icon (Option C)

Use lucide-react's `Trophy` icon with custom styling.

**Rationale:** More visually appealing than a simple badge.

### 4. Display Choices

- **Team Name**: Full name (e.g., "Team Liquid", not "TL")
- **Player Count**: Remove (always 5 in Dota 2)
- **Team Logo**: Don't display (already shown above in SeriesHeader)
- **Score Format**: Number only (e.g., `42`)

## Current State

### File: `ui/src/components/series/PlayersTable.tsx`

**Current Header (Lines 44-46):**

```tsx
<div className={cn("text-sm font-semibold mb-3 flex items-center gap-2")}>
  <span>{team}</span>
  <span className="text-foreground-muted">({players.length})</span>
</div>
```

**Current Props:**

```tsx
interface PlayersTableProps {
  match: SeriesMatchDetail;
}

interface TeamPlayersTableProps {
  players: PlayerData[];
  team: "Radiant" | "Dire";
  captain: number | null;
  teamColor: "primary" | "secondary";
}
```

## Available Data

### From `SeriesMatchDetail` (already in `match` prop):

- `radiant_win: boolean` - which team won
- `radiant_score: number` - radiant team's score
- `dire_score: number` - dire team's score

### From `SeriesDetail` (in parent component):

- `team_a: TeamInfo` - { id, name, tag, logo_url }
- `team_b: TeamInfo` - { id, name, tag, logo_url }

### Type Definition:

```tsx
export type TeamInfo = {
  id: number;
  name: string;
  tag?: string;
  logo_url?: string;
};
```

## Theme Colors

- **Primary (Dire Red)**: `hsl(0, 84%, 50%)` → `text-primary-500`
- **Secondary (Radiant Gold)**: `hsl(38, 92%, 50%)` → `text-secondary-500`
- **Success (Green)**: `hsl(142, 71%, 45%)` → `bg-success-500`

## Implementation Steps

### Step 1: Create `TeamIcon` Component

**File:** `ui/src/components/series/TeamIcon.tsx` (NEW)

```tsx
import { cn } from "@/lib";

interface TeamIconProps {
  type: "radiant" | "dire";
  className?: string;
}

export function TeamIcon({ type, className }: TeamIconProps) {
  // User will provide SVG paths for radiant and dire
  const radiantPath = "..."; // TODO: Add radiant SVG path
  const direPath = "..."; // TODO: Add dire SVG path

  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      className={cn(
        "w-5 h-5",
        type === "radiant" ? "text-secondary-500" : "text-primary-500",
        className,
      )}
    >
      <path
        d={type === "radiant" ? radiantPath : direPath}
        fill="currentColor"
      />
    </svg>
  );
}
```

### Step 2: Update `PlayersTable.tsx` Props

```tsx
import type { SeriesMatchDetail, PlayerData, TeamInfo } from "@/api/api";

interface PlayersTableProps {
  match: SeriesMatchDetail;
  radiantTeam: TeamInfo;
  direTeam: TeamInfo;
}

interface TeamPlayersTableProps {
  players: PlayerData[];
  team: "Radiant" | "Dire";
  captain: number | null;
  teamColor: "primary" | "secondary";
  teamInfo: TeamInfo;
  score: number;
  isWinner: boolean;
}
```

### Step 3: Update `PlayersTable` Component

```tsx
export function PlayersTable({
  match,
  radiantTeam,
  direTeam,
}: PlayersTableProps) {
  const {
    players_data,
    radiant_captain,
    dire_captain,
    radiant_win,
    radiant_score,
    dire_score,
  } = match;

  const radiantPlayers = players_data.filter((p) => p.is_radiant);
  const direPlayers = players_data.filter((p) => !p.is_radiant);

  return (
    <div className="flex flex-col gap-3">
      <TeamPlayersTable
        players={radiantPlayers}
        team="Radiant"
        captain={radiant_captain}
        teamColor="secondary"
        teamInfo={radiantTeam}
        score={radiant_score}
        isWinner={radiant_win}
      />
      <TeamPlayersTable
        players={direPlayers}
        team="Dire"
        captain={dire_captain}
        teamColor="primary"
        teamInfo={direTeam}
        score={dire_score}
        isWinner={!radiant_win}
      />
    </div>
  );
}
```

### Step 4: Update `TeamPlayersTable` Header

```tsx
import { Trophy } from "lucide-react";
import { TeamIcon } from "./TeamIcon";

function TeamPlayersTable({
  players,
  team,
  captain,
  teamColor,
  teamInfo,
  score,
  isWinner,
}: TeamPlayersTableProps) {
  return (
    <div className="bg-background-accent rounded-lg p-4">
      {/* New Team Header */}
      <div className="flex items-center gap-2 mb-3">
        <TeamIcon type={team.toLowerCase() as "radiant" | "dire"} />
        <span className="text-sm font-semibold">{teamInfo.name}</span>
        <span className="text-foreground-muted text-sm">{score}</span>
        {isWinner && (
          <div className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs bg-success-500 text-white font-semibold">
            <Trophy className="w-3 h-3" />
            WON
          </div>
        )}
      </div>

      {/* Rest of the table remains unchanged */}
      <Table className="table-fixed">{/* ... existing table code ... */}</Table>
    </div>
  );
}
```

### Step 5: Update `GameTabContent.tsx`

```tsx
import type { SeriesMatchDetail, TeamInfo } from "@/api/api";

interface GameTabContentProps {
  match: SeriesMatchDetail;
  gameNumber: number;
  radiantTeam: TeamInfo;
  direTeam: TeamInfo;
}

export function GameTabContent({
  match,
  radiantTeam,
  direTeam,
}: GameTabContentProps) {
  // ... existing code ...

  return (
    <div className="rounded-xl border border-white/10 bg-background-card focus:outline-none focus-visible:ring-0 overflow-hidden">
      {/* ... existing header ... */}
      <div className="px-6 pb-6">
        <PlayersTable
          match={match}
          radiantTeam={radiantTeam}
          direTeam={direTeam}
        />
      </div>
    </div>
  );
}
```

### Step 6: Update `$id.tsx` Route

```tsx
{
  data.matches.map((match, i) => (
    <TabsContent key={i} value={`game-${i + 1}`}>
      <GameTabContent
        match={match}
        gameNumber={i + 1}
        radiantTeam={data.team_a}
        direTeam={data.team_b}
      />
    </TabsContent>
  ));
}
```

### Step 7: Update `index.ts` Export

**File:** `ui/src/components/series/index.ts`

```tsx
export { PlayersTable } from "./PlayersTable";
export { TeamIcon } from "./TeamIcon"; // ADD THIS
```

## Visual Layout

```
┌─────────────────────────────────────────────────────────────┐
│  [☀️] Team Liquid  42  [🏆 WON]                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Hero │ Player │ Items │ NET │ K/D/A │ LH            │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ ...                                                 │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│  [🔥] Team Secret  38                                      │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Hero │ Player │ Items │ NET │ K/D/A │ LH            │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ ...                                                 │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Files to Modify

1. **NEW:** `ui/src/components/series/TeamIcon.tsx` - Team icon component
2. **MODIFY:** `ui/src/components/series/PlayersTable.tsx` - Main changes
3. **MODIFY:** `ui/src/components/series/GameTabContent.tsx` - Add team props
4. **MODIFY:** `ui/src/routes/series/$id.tsx` - Pass team info
5. **MODIFY:** `ui/src/components/series/index.ts` - Export TeamIcon

## Notes

- User needs to provide SVG paths for Radiant and Dire icons
- The `Trophy` icon is already available from lucide-react (already imported in PlayersTable.tsx)
- No API changes required
- Team logos are not displayed (already shown in SeriesHeader above)
- Player count is removed (always 5 in Dota 2)
