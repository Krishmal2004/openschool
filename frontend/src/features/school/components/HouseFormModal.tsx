import { TextInput } from "@carbon/react";
import FormModal from "@/shared/ui/FormModal";

export interface HouseForm {
  name: string;
  code: string;
  color: string;
}

interface Props {
  mode: "create" | "edit";
  form: HouseForm;
  onChange: (form: HouseForm) => void;
  nameTouched: boolean;
  onNameBlur: () => void;
  isPending: boolean;
  isError: boolean;
  error: unknown;
  onClose: () => void;
  onSubmit: () => void;
}

export default function HouseFormModal({ mode, form, onChange, nameTouched, onNameBlur, isPending, isError, error, onClose, onSubmit }: Props) {
  return (
    <FormModal
      open
      title={mode === "create" ? "Add house" : "Edit house"}
      onClose={onClose}
      onSubmit={onSubmit}
      isPending={isPending}
      submitDisabled={!form.name.trim()}
      isError={isError}
      error={error}
      errorFallback="Failed to save house"
    >
      <TextInput
        id="house-name"
        labelText="Name"
        placeholder="e.g. Vijaya"
        value={form.name}
        onChange={(e) => onChange({ ...form, name: e.target.value })}
        onBlur={onNameBlur}
        invalid={nameTouched && !form.name.trim()}
        invalidText="House name is required."
        className="os-mb-4"
      />
      <TextInput id="house-code" labelText="Short code (optional)" placeholder="e.g. VJ" value={form.code} onChange={(e) => onChange({ ...form, code: e.target.value })} className="os-mb-4" />
      <label htmlFor="house-color" className="os-text-xs os-block os-mb-2">Color</label>
      <div className="os-flex os-items-center os-gap-3">
        <input id="house-color" type="color" value={form.color} onChange={(e) => onChange({ ...form, color: e.target.value })} className="os-w-3 os-h-2h os-p-0 os-border-tertiary os-pointer" />
        <TextInput id="house-color-hex" labelText="" hideLabel value={form.color} onChange={(e) => onChange({ ...form, color: e.target.value })} className="os-max-w-8" />
      </div>
    </FormModal>
  );
}
