import { Button, TextInput, DatePicker, DatePickerInput } from "@carbon/react";
import { Checkmark, Add } from "@carbon/icons-react";
import type { Term } from "@/features/school/api/term";
import { toYmd, isDateRangeInvalid } from "@/shared/lib/date";
import type { TermFormValues, TermTouched } from "@/features/school/lib/termForm";

interface Props {
  form: TermFormValues;
  onChange: (form: TermFormValues) => void;
  touched: TermTouched;
  onTouch: (field: keyof TermFormValues) => void;
  editing: Term | null;
  isSaving: boolean;
  onSubmit: () => void;
  onCancelEdit: () => void;
}

// Inline add / edit form under the term list.
export default function TermForm({ form, onChange, touched, onTouch, editing, isSaving, onSubmit, onCancelEdit }: Props) {
  const rangeInvalid = isDateRangeInvalid(form.start_date, form.end_date);
  // Remounting the pickers on term switch makes flatpickr pick up the prefilled value.
  const pickerKey = editing?.id ?? "new";
  return (
    <div className="os-grid os-gap-3">
      {editing && <p className="os-m-0 os-text-xs os-fw-600 os-c-primary">Editing {editing.name}</p>}
      <TextInput
        id="term-name"
        labelText="Term name"
        placeholder="e.g. Term 1"
        value={form.name}
        onChange={(e) => onChange({ ...form, name: e.target.value })}
        onBlur={() => onTouch("name")}
        invalid={!!touched.name && !form.name.trim()}
        invalidText="A name is required."
      />
      <DatePicker key={`start-${pickerKey}`} datePickerType="single" dateFormat="Y-m-d" value={form.start_date} onChange={(dates) => onChange({ ...form, start_date: toYmd(dates[0]) })}>
        <DatePickerInput id="term-start" labelText="Start Date" placeholder="YYYY-MM-DD" onBlur={() => onTouch("start_date")} invalid={!!touched.start_date && !form.start_date} invalidText="A start date is required." />
      </DatePicker>
      <DatePicker key={`end-${pickerKey}`} datePickerType="single" dateFormat="Y-m-d" value={form.end_date} onChange={(dates) => onChange({ ...form, end_date: toYmd(dates[0]) })}>
        <DatePickerInput
          id="term-end"
          labelText="End Date"
          placeholder="YYYY-MM-DD"
          onBlur={() => onTouch("end_date")}
          invalid={!!touched.end_date && (!form.end_date || rangeInvalid)}
          invalidText={rangeInvalid ? "End date must be after the start date." : "An end date is required."}
        />
      </DatePicker>
      <div className="os-flex os-gap-2">
        <Button kind="ghost" size="sm" renderIcon={editing ? Checkmark : Add} onClick={onSubmit} disabled={isSaving}>
          {isSaving ? (editing ? "Saving…" : "Adding…") : editing ? "Save Changes" : "Add Term"}
        </Button>
        {editing && <Button kind="ghost" size="sm" onClick={onCancelEdit} disabled={isSaving}>Cancel</Button>}
      </div>
    </div>
  );
}
