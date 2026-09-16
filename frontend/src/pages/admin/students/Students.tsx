import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Add, Edit, TrashCan, Close } from "@carbon/icons-react";
import {
  Button,
  IconButton,
  Select,
  SelectItem,
  Pagination,
  Tag,
  DataTable,
  TableContainer,
  TableToolbar,
  TableToolbarContent,
  TableToolbarSearch,
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
} from "@carbon/react";
import EntityCombobox from "../../../components/common/EntityCombobox";
import { useStudents, useDeleteStudent } from "../../../queries/useStudents";
import { useGrades } from "../../../queries/useGrades";
import { useHouses } from "../../../queries/useHouses";
import { useCurrentClasses } from "../../../queries/useClasses";
import type { Student } from "../../../services/student";
import { usePagination } from "../../../hooks/usePagination";
import TableSkeleton from "../../../components/common/TableSkeleton";
import ErrorMessage from "../../../components/common/ErrorMessage";
import EmptyState from "../../../components/common/EmptyState";
import ConfirmDeleteModal from "../../../components/common/ConfirmDeleteModal";
import AgentFindingsBanner from "../../../components/common/AgentFindingsBanner";
import MutationErrorNotification from "../../../components/common/MutationErrorNotification";

const STUDENT_TABLE_HEADERS = [
  { key: "index_number", header: "Index No." },
  { key: "full_name", header: "Full Name" },
  { key: "class_name", header: "Class" },
  { key: "house_name", header: "House" },
  { key: "phone", header: "Phone" },
  { key: "whatsapp", header: "WhatsApp" },
  { key: "enrollment_status", header: "Status" },
  { key: "actions", header: "Actions" },
];

export default function Students() {
  const navigate = useNavigate();
  const { data: students, isLoading, isError, refetch } = useStudents();
  const { data: grades } = useGrades();
  const { data: houses } = useHouses();
  const { data: classes } = useCurrentClasses();
  const deleteStudent = useDeleteStudent();
  const [query, setQuery] = useState("");
  const [grade, setGrade] = useState("");
  const [cls, setCls] = useState("");
  const [gender, setGender] = useState("");
  const [house, setHouse] = useState("");
  const [toDelete, setToDelete] = useState<Student | null>(null);

  const classOptions = (classes ?? []).filter(
    (c) => !grade || c.grade_name === grade,
  );

  const changeGrade = (value: string) => {
    setGrade(value);
    setCls("");
  };

  const clearAllFilters = () => {
    setQuery("");
    setGrade("");
    setCls("");
    setGender("");
    setHouse("");
  };

  const hasActiveFilters = Boolean(query || grade || cls || gender || house);

  const filtered = (students ?? []).filter((s) => {
    const q = query.toLowerCase();
    const matchesSearch =
      s.full_name.toLowerCase().includes(q) ||
      s.index_number.toLowerCase().includes(q);
    const matchesGrade = !grade || s.grade_name === grade;
    const matchesClass = !cls || s.class_name === cls;
    const matchesGender = !gender || s.gender === gender;
    const matchesHouse = !house || s.house_name === house;
    return (
      matchesSearch &&
      matchesGrade &&
      matchesClass &&
      matchesGender &&
      matchesHouse
    );
  });

  const { page, pageSize, pageItems, totalItems, onChange } = usePagination(filtered, 10);

  const handleDelete = () => {
    if (!toDelete) return;
    deleteStudent.mutate(toDelete.id, { onSettled: () => setToDelete(null) });
  };

  return (
    <div className="os-page">
      <div className="os-page__header">
        <div className="os-page__header-left">
          <h1 className="os-page__title">Students</h1>
          <p className="os-page__subtitle">Manage student enrolment and profiles</p>
        </div>
        <Button renderIcon={Add} kind="primary" size="md" as={Link} to="/students/new">
          Enrol Student
        </Button>
      </div>

      <AgentFindingsBanner
        titles={[
          "Students with no guardian on file",
          "Student gender / school-type mismatches",
          "Students with no current-year class",
          "Student accounts stuck in first-login setup",
        ]}
      />

      <div className="os-section" style={{ background: "#ffffff" }}>
        <div className="os-toolbar" style={{ padding: "1rem 1.5rem 0.5rem" }}>
          <div style={{ flex: "1 1 16rem", minWidth: "12rem" }}>
            <TableToolbarSearch
              persistent
              placeholder="Search by name or index number…"
              value={query}
              onChange={(e) => setQuery(typeof e === "string" ? e : e.target.value)}
            />
          </div>
          <div style={{ minWidth: "9rem" }}>
            <EntityCombobox
              id="filter-grade"
              items={grades ?? []}
              selectedId={grade}
              onSelect={changeGrade}
              getId={(g) => g.name}
              itemToString={(g) => g.name}
              placeholder="All grades"
            />
          </div>
          <div style={{ minWidth: "9rem" }}>
            <EntityCombobox
              id="filter-class"
              items={classOptions}
              selectedId={cls}
              onSelect={setCls}
              getId={(c) => c.name}
              itemToString={(c) => c.name}
              placeholder="All classes"
            />
          </div>
          <div style={{ minWidth: "8rem" }}>
            <Select
              id="filter-gender"
              labelText=""
              size="md"
              value={gender}
              onChange={(e) => setGender(e.target.value)}
            >
              <SelectItem value="" text="Any gender" />
              <SelectItem value="male" text="Male" />
              <SelectItem value="female" text="Female" />
            </Select>
          </div>
          <div style={{ minWidth: "9rem" }}>
            <Select
              id="filter-house"
              labelText=""
              size="md"
              value={house}
              onChange={(e) => setHouse(e.target.value)}
            >
              <SelectItem value="" text="All houses" />
              {houses?.map((h) => (
                <SelectItem key={h.id} value={h.name} text={h.name} />
              ))}
            </Select>
          </div>
        </div>

        {hasActiveFilters && (
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", flexWrap: "wrap", padding: "0 1.5rem 0.75rem" }}>
            <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)" }}>Active Filters:</span>
            {grade && <Tag type="cyan" filter onClose={() => changeGrade("")}>Grade: {grade}</Tag>}
            {cls && <Tag type="cyan" filter onClose={() => setCls("")}>Class: {cls}</Tag>}
            {gender && <Tag type="teal" filter onClose={() => setGender("")}>Gender: {gender}</Tag>}
            {house && <Tag type="purple" filter onClose={() => setHouse("")}>House: {house}</Tag>}
            {query && <Tag type="blue" filter onClose={() => setQuery("")}>Search: "{query}"</Tag>}
            <Button kind="ghost" size="sm" renderIcon={Close} onClick={clearAllFilters}>
              Clear All
            </Button>
          </div>
        )}

        <MutationErrorNotification
          isError={deleteStudent.isError}
          error={deleteStudent.error}
          title="Could not delete student"
          fallback="Failed to delete student."
          onClose={() => deleteStudent.reset()}
          style={{ margin: "0 1.5rem 1rem" }}
        />

        {isLoading ? (
          <TableSkeleton headers={STUDENT_TABLE_HEADERS.map(h => h.header)} />
        ) : isError ? (
          <ErrorMessage message="Failed to load students" onRetry={refetch} />
        ) : filtered.length === 0 ? (
          <EmptyState
            title="No students found"
            description="Enrol your first student or adjust your search filters to get started."
          />
        ) : (
          <>
            <DataTable
              rows={pageItems.map((s) => ({ ...s, id: s.id }))}
              headers={STUDENT_TABLE_HEADERS}
              render={({
                rows,
                headers,
                getHeaderProps,
                getRowProps,
              }) => (
                <TableContainer className="os-table-container">
                  <TableToolbar style={{ background: "#ffffff" }}>
                    <TableToolbarContent>
                      <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--os-text-tertiary)", paddingRight: "1rem" }}>
                        Showing {pageItems.length} of {filtered.length} students
                      </span>
                    </TableToolbarContent>
                  </TableToolbar>

                  <Table className="os-table">
                    <TableHead>
                      <TableRow>
                        {headers.map((header) => (
                          <TableHeader {...getHeaderProps({ header })}>
                            {header.header}
                          </TableHeader>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {rows.map((row) => {
                        const student = pageItems.find((s) => s.id === row.id);
                        if (!student) return null;
                        return (
                          <TableRow {...getRowProps({ row })}>
                            <TableCell>
                              <span className="os-table__mono">{student.index_number}</span>
                            </TableCell>
                            <TableCell>
                              <Link to={`/students/${student.id}`} className="os-table__link">
                                {student.full_name}
                              </Link>
                            </TableCell>
                            <TableCell className="os-table__muted">{student.class_name ?? "-"}</TableCell>
                            <TableCell className="os-table__muted">{student.house_name ?? "-"}</TableCell>
                            <TableCell className="os-table__muted">{student.phone ?? "-"}</TableCell>
                            <TableCell className="os-table__muted">{student.whatsapp ?? "-"}</TableCell>
                            <TableCell>
                              {student.enrollment_status === "active" ? (
                                <Tag type="green" size="sm">Active</Tag>
                              ) : (
                                <Tag type="red" size="sm">Left</Tag>
                              )}
                            </TableCell>
                            <TableCell>
                              <div style={{ display: "flex", gap: "0.25rem", justifyContent: "flex-end" }}>
                                <IconButton
                                  label="Edit profile"
                                  kind="ghost"
                                  size="sm"
                                  onClick={() =>
                                    navigate(`/students/${student.id}`, {
                                      state: { edit: true },
                                    })
                                  }
                                >
                                  <Edit />
                                </IconButton>
                                <IconButton
                                  label="Delete student"
                                  kind="ghost"
                                  size="sm"
                                  onClick={() => setToDelete(student)}
                                >
                                  <TrashCan />
                                </IconButton>
                              </div>
                            </TableCell>
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                </TableContainer>
              )}
            />
            <Pagination
              totalItems={totalItems}
              page={page}
              pageSize={pageSize}
              pageSizes={[10, 20, 50]}
              onChange={onChange}
            />
          </>
        )}
      </div>

      <ConfirmDeleteModal
        open={!!toDelete}
        title="Delete student"
        description={
          <>
            Delete <strong>{toDelete?.full_name}</strong>? This removes their account and cannot be undone.
          </>
        }
        isPending={deleteStudent.isPending}
        onClose={() => setToDelete(null)}
        onConfirm={handleDelete}
      />
    </div>
  );
}


