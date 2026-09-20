import type { Dispatch, SetStateAction } from "react";
import { Button, TextInput } from "@carbon/react";
import { Home, Add } from "@carbon/icons-react";
import StepShell from "@/features/school/components/setup/StepShell";
import RepeatableRow from "@/features/school/components/setup/RepeatableRow";
import { HOUSE_COLOR_PALETTE } from "@/features/school/setupConstants";

export interface HouseRow {
  name: string;
  code: string;
  color: string;
}

interface Props {
  houses: HouseRow[];
  setHouses: Dispatch<SetStateAction<HouseRow[]>>;
}

export default function HousesStep({ houses, setHouses }: Props) {
  return (
    <StepShell icon={Home} title="Houses" subtitle="Optional - students and staff are auto-assigned to whichever house has the fewest members.">
      {houses.map((h, i) => (
        <RepeatableRow key={i} onRemove={() => setHouses((hs) => hs.filter((_, idx) => idx !== i))}>
          <TextInput
            id={`house-name-${i}`}
            labelText="Name"
            placeholder="e.g. Vijaya"
            size="md"
            value={h.name}
            onChange={(e) =>
              setHouses((hs) => hs.map((row, idx) => (idx === i ? { ...row, name: e.target.value } : row)))
            }
          />
          <TextInput
            id={`house-code-${i}`}
            labelText="Code (optional)"
            placeholder="e.g. VJ"
            size="md"
            value={h.code}
            onChange={(e) =>
              setHouses((hs) => hs.map((row, idx) => (idx === i ? { ...row, code: e.target.value } : row)))
            }
          />
          <div>
            <label htmlFor={`house-color-${i}`} className="os-text-xs os-block os-mb-1">
              Colour
            </label>
            <input
              id={`house-color-${i}`}
              type="color"
              value={h.color}
              onChange={(e) =>
                setHouses((hs) => hs.map((row, idx) => (idx === i ? { ...row, color: e.target.value } : row)))
              } className="os-w-2h os-h-2h os-p-0 os-border-tertiary os-pointer"
            />
          </div>
        </RepeatableRow>
      ))}
      <Button
        kind="ghost"
        size="sm"
        renderIcon={Add}
        onClick={() =>
          setHouses((hs) => [...hs, { name: "", code: "", color: HOUSE_COLOR_PALETTE[hs.length % HOUSE_COLOR_PALETTE.length] }])
        }
      >
        Add another house
      </Button>
    </StepShell>
  );
}
