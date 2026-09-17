import { ComboBox } from "@carbon/react";

interface EntityComboboxProps<T> {
  id: string;
  items: T[];
  selectedId: string;
  onSelect: (id: string) => void;
  itemToString: (item: T) => string;
  getId: (item: T) => string;
  labelText?: string;
  placeholder?: string;
  invalid?: boolean;
  invalidText?: string;
  disabled?: boolean;
}

// Type-to-filter picker for large people lists where a plain Select would not scale.
export default function EntityCombobox<T>({
  id,
  items,
  selectedId,
  onSelect,
  itemToString,
  getId,
  labelText,
  placeholder,
  invalid,
  invalidText,
  disabled,
}: EntityComboboxProps<T>) {
  const selectedItem = items.find((item) => getId(item) === selectedId) ?? null;

  return (
    <ComboBox
      id={id}
      items={items}
      itemToString={(item) => (item ? itemToString(item as T) : "")}
      selectedItem={selectedItem}
      onChange={({ selectedItem }) => onSelect(selectedItem ? getId(selectedItem) : "")}
      titleText={labelText}
      placeholder={placeholder ?? "Search by name or ID…"}
      invalid={invalid}
      invalidText={invalidText}
      disabled={disabled}
    />
  );
}
