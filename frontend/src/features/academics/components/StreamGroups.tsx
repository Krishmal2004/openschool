import { useState } from "react";
import { Add } from "@carbon/icons-react";
import { Button, SkeletonText, Tag, TextInput } from "@carbon/react";
import { useStreamGroups, useCreateStreamGroup } from "@/features/academics/queries/useClasses";
import type { Stream } from "@/features/academics/api/stream";

// Sub-groups under one A/L stream, with an inline add box.
export default function StreamGroups({ stream }: { stream: Stream }) {
  const { data: groups, isLoading } = useStreamGroups(stream.id);
  const createGroup = useCreateStreamGroup();
  const [name, setName] = useState("");

  const add = () => {
    if (!name.trim()) return;
    createGroup.mutate({ streamId: stream.id, data: { name: name.trim() } }, { onSuccess: () => setName("") });
  };

  return (
    <div className="os-pl-6">
      <div className="os-flex os-wrap os-gap-2 os-mb-2">
        {isLoading ? (
          <SkeletonText width="30%" />
        ) : groups?.length ? (
          groups.map((g) => <Tag key={g.id} type="teal" size="sm">{g.name}</Tag>)
        ) : (
          <span className="os-text-xs os-c-tertiary">No sub-groups</span>
        )}
      </div>
      <div className="os-flex os-gap-2 os-items-end">
        <TextInput id={`new-group-${stream.id}`} labelText="New group name" hideLabel placeholder="e.g. Physical Science" size="sm" value={name} onChange={(e) => setName(e.target.value)} className="os-max-w-14" />
        <Button kind="ghost" size="sm" renderIcon={Add} disabled={!name.trim() || createGroup.isPending} onClick={add}>Add group</Button>
      </div>
    </div>
  );
}
