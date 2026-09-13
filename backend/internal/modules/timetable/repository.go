package timetable

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type classroomRepository struct{ queries *db.Queries }

func newClassroomRepository(pool *pgxpool.Pool) *classroomRepository {
	return &classroomRepository{queries: db.New(pool)}
}
func (r *classroomRepository) create(ctx context.Context, command classroomCommand) (Classroom, error) {
	row, err := r.queries.CreateClassroom(ctx, db.CreateClassroomParams{
		Name: command.Name, Code: optionalClassroomText(command.Code), Capacity: optionalClassroomInt(command.Capacity),
		RoomType: command.RoomType, SubjectID: optionalClassroomUUID(command.SubjectID),
	})
	return mapClassroom(row), err
}
func (r *classroomRepository) list(ctx context.Context) ([]ClassroomListItem, error) {
	rows, err := r.queries.ListClassrooms(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]ClassroomListItem, len(rows))
	for i, row := range rows {
		result[i] = ClassroomListItem{Classroom: Classroom{
			ID: row.ID, Name: row.Name, Code: classroomText(row.Code), Capacity: classroomInt(row.Capacity),
			CreatedAt: row.CreatedAt.Time, RoomType: row.RoomType, SubjectID: classroomUUID(row.SubjectID),
		}, SubjectName: classroomText(row.SubjectName)}
	}
	return result, nil
}
func (r *classroomRepository) update(ctx context.Context, id uuid.UUID, command classroomCommand) (Classroom, error) {
	row, err := r.queries.UpdateClassroom(ctx, db.UpdateClassroomParams{
		ID: id, Name: command.Name, Code: optionalClassroomText(command.Code), Capacity: optionalClassroomInt(command.Capacity),
		RoomType: command.RoomType, SubjectID: optionalClassroomUUID(command.SubjectID),
	})
	return mapClassroom(row), err
}
func (r *classroomRepository) delete(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteClassroom(ctx, id)
}
func optionalClassroomText(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}
func optionalClassroomInt(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}
func optionalClassroomUUID(value *uuid.UUID) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *value, Valid: true}
}
func classroomText(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
func classroomInt(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	result := value.Int32
	return &result
}
func classroomUUID(value pgtype.UUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}
	result := uuid.UUID(value.Bytes)
	return &result
}
func mapClassroom(row db.Classroom) Classroom {
	return Classroom{ID: row.ID, Name: row.Name, Code: classroomText(row.Code), Capacity: classroomInt(row.Capacity), CreatedAt: row.CreatedAt.Time, RoomType: row.RoomType, SubjectID: classroomUUID(row.SubjectID)}
}

type subjectPeriodRequirementRepository struct{ queries *db.Queries }

func newSubjectPeriodRequirementRepository(pool *pgxpool.Pool) *subjectPeriodRequirementRepository {
	return &subjectPeriodRequirementRepository{queries: db.New(pool)}
}
func (r *subjectPeriodRequirementRepository) upsert(ctx context.Context, command subjectPeriodRequirementCommand) (SubjectPeriodRequirement, error) {
	row, err := r.queries.UpsertSubjectPeriodRequirement(ctx, db.UpsertSubjectPeriodRequirementParams{
		AcademicYearID: command.AcademicYearID, GradeID: command.GradeID, SubjectID: command.SubjectID,
		PeriodsPerWeek: command.PeriodsPerWeek, LabPeriodsPerWeek: command.LabPeriodsPerWeek,
		DoublePeriodBlocks: command.DoublePeriodBlocks,
	})
	return mapSubjectPeriodRequirement(row), err
}
func (r *subjectPeriodRequirementRepository) listByGrade(ctx context.Context, academicYearID, gradeID uuid.UUID) ([]SubjectPeriodRequirementListItem, error) {
	rows, err := r.queries.ListSubjectPeriodRequirementsByGrade(ctx, db.ListSubjectPeriodRequirementsByGradeParams{AcademicYearID: academicYearID, GradeID: gradeID})
	if err != nil {
		return nil, err
	}
	result := make([]SubjectPeriodRequirementListItem, len(rows))
	for i, row := range rows {
		result[i] = SubjectPeriodRequirementListItem{SubjectPeriodRequirement: SubjectPeriodRequirement{
			ID: row.ID, AcademicYearID: row.AcademicYearID, GradeID: row.GradeID, SubjectID: row.SubjectID,
			PeriodsPerWeek: row.PeriodsPerWeek, CreatedAt: row.CreatedAt.Time,
			LabPeriodsPerWeek: row.LabPeriodsPerWeek, DoublePeriodBlocks: row.DoublePeriodBlocks,
		}, SubjectName: row.SubjectName, SubjectCode: row.SubjectCode}
	}
	return result, nil
}
func (r *subjectPeriodRequirementRepository) deleteRequirement(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteSubjectPeriodRequirement(ctx, id)
}
func mapSubjectPeriodRequirement(row db.SubjectPeriodRequirement) SubjectPeriodRequirement {
	return SubjectPeriodRequirement{
		ID: row.ID, AcademicYearID: row.AcademicYearID, GradeID: row.GradeID, SubjectID: row.SubjectID,
		PeriodsPerWeek: row.PeriodsPerWeek, CreatedAt: row.CreatedAt.Time,
		LabPeriodsPerWeek: row.LabPeriodsPerWeek, DoublePeriodBlocks: row.DoublePeriodBlocks,
	}
}

type teacherAvailabilityRepository struct{ queries *db.Queries }

func newTeacherAvailabilityRepository(pool *pgxpool.Pool) *teacherAvailabilityRepository {
	return &teacherAvailabilityRepository{queries: db.New(pool)}
}
func (r *teacherAvailabilityRepository) createAvailability(ctx context.Context, teacherID uuid.UUID, command teacherAvailabilityCommand) (TeacherAvailability, error) {
	row, err := r.queries.CreateTeacherAvailability(ctx, db.CreateTeacherAvailabilityParams{
		TeacherID: teacherID, AcademicYearID: command.AcademicYearID,
		DayOfWeek: command.DayOfWeek, PeriodNumber: command.PeriodNumber,
	})
	return mapTeacherAvailability(row), err
}
func (r *teacherAvailabilityRepository) listAvailability(ctx context.Context, teacherID, academicYearID uuid.UUID) ([]TeacherAvailability, error) {
	rows, err := r.queries.ListTeacherAvailabilityByTeacherYear(ctx, db.ListTeacherAvailabilityByTeacherYearParams{TeacherID: teacherID, AcademicYearID: academicYearID})
	if err != nil {
		return nil, err
	}
	result := make([]TeacherAvailability, len(rows))
	for i, row := range rows {
		result[i] = mapTeacherAvailability(row)
	}
	return result, nil
}
func (r *teacherAvailabilityRepository) deleteAvailability(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteTeacherAvailability(ctx, id)
}
func mapTeacherAvailability(row db.TeacherAvailability) TeacherAvailability {
	return TeacherAvailability{
		ID: row.ID, TeacherID: row.TeacherID, AcademicYearID: row.AcademicYearID,
		DayOfWeek: row.DayOfWeek, PeriodNumber: row.PeriodNumber, CreatedAt: row.CreatedAt.Time,
	}
}

type gradeSectionRepository struct{ queries *db.Queries }

func newGradeSectionRepository(pool *pgxpool.Pool) *gradeSectionRepository {
	return &gradeSectionRepository{queries: db.New(pool)}
}
func (r *gradeSectionRepository) createSection(ctx context.Context, command gradeSectionCommand, start, end int64) (gradeSectionRecord, error) {
	row, err := r.queries.CreateGradeSection(ctx, db.CreateGradeSectionParams{
		AcademicYearID: command.AcademicYearID, Name: command.Name,
		IntervalStartTime: pgtype.Time{Microseconds: start, Valid: true}, IntervalEndTime: pgtype.Time{Microseconds: end, Valid: true},
		SectionHeadTeacherID: optionalClassroomUUID(command.SectionHeadTeacherID), SortOrder: command.SortOrder,
	})
	return mapGradeSectionRecord(row), err
}
func (r *gradeSectionRepository) getSection(ctx context.Context, id uuid.UUID) (gradeSectionRecord, error) {
	row, err := r.queries.GetGradeSectionByID(ctx, id)
	return mapGradeSectionRecord(row), err
}
func (r *gradeSectionRepository) listSections(ctx context.Context, yearID uuid.UUID) ([]gradeSectionRecord, error) {
	rows, err := r.queries.ListGradeSectionsByYear(ctx, yearID)
	if err != nil {
		return nil, err
	}
	result := make([]gradeSectionRecord, len(rows))
	for i, row := range rows {
		result[i] = gradeSectionRecord{
			ID: row.ID, AcademicYearID: row.AcademicYearID, Name: row.Name,
			IntervalStartMicroseconds: row.IntervalStartTime.Microseconds, IntervalEndMicroseconds: row.IntervalEndTime.Microseconds,
			SectionHeadTeacherID: classroomUUID(row.SectionHeadTeacherID), SectionHeadName: classroomText(row.SectionHeadName), SortOrder: row.SortOrder,
		}
	}
	return result, nil
}
func (r *gradeSectionRepository) updateSection(ctx context.Context, id uuid.UUID, command updateGradeSectionCommand, start, end int64) error {
	_, err := r.queries.UpdateGradeSection(ctx, db.UpdateGradeSectionParams{
		ID: id, Name: command.Name, IntervalStartTime: pgtype.Time{Microseconds: start, Valid: true},
		IntervalEndTime:      pgtype.Time{Microseconds: end, Valid: true},
		SectionHeadTeacherID: optionalClassroomUUID(command.SectionHeadTeacherID), SortOrder: command.SortOrder,
	})
	return err
}
func (r *gradeSectionRepository) deleteSection(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.queries.DeleteGradeSection(ctx, id)
}
func (r *gradeSectionRepository) assignGrade(ctx context.Context, sectionID, gradeID, yearID uuid.UUID) error {
	return r.queries.AssignGradeToSection(ctx, db.AssignGradeToSectionParams{GradeSectionID: sectionID, GradeID: gradeID, AcademicYearID: yearID})
}
func (r *gradeSectionRepository) removeGrade(ctx context.Context, sectionID, gradeID uuid.UUID) error {
	return r.queries.RemoveGradeFromSection(ctx, db.RemoveGradeFromSectionParams{GradeSectionID: sectionID, GradeID: gradeID})
}
func (r *gradeSectionRepository) listGradeIDs(ctx context.Context, sectionID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.queries.ListGradesBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	result := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		result[i] = row.ID
	}
	return result, nil
}
func (r *gradeSectionRepository) createPeriod(ctx context.Context, values periodValues) (TimetablePeriod, error) {
	row, err := r.queries.CreateTimetablePeriod(ctx, db.CreateTimetablePeriodParams{
		GradeSectionID: values.GradeSectionID, SortOrder: values.SortOrder, PeriodNumber: optionalClassroomInt(values.PeriodNumber),
		StartTime: pgtype.Time{Microseconds: values.StartMicroseconds, Valid: true}, EndTime: pgtype.Time{Microseconds: values.EndMicroseconds, Valid: true}, SlotType: values.SlotType,
	})
	return mapTimetablePeriod(row), err
}
func (r *gradeSectionRepository) listPeriods(ctx context.Context, sectionID uuid.UUID) ([]TimetablePeriod, error) {
	rows, err := r.queries.ListTimetablePeriodsBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	result := make([]TimetablePeriod, len(rows))
	for i, row := range rows {
		result[i] = mapTimetablePeriod(row)
	}
	return result, nil
}
func (r *gradeSectionRepository) deletePeriods(ctx context.Context, sectionID uuid.UUID) error {
	return r.queries.DeleteTimetablePeriodsBySection(ctx, sectionID)
}
func (r *gradeSectionRepository) getSettingsValues(ctx context.Context, yearID uuid.UUID) (settingsValues, error) {
	row, err := r.queries.GetTimetableSettingsByYear(ctx, yearID)
	if err != nil {
		return settingsValues{}, err
	}
	return settingsValues{
		AcademicYearID: row.AcademicYearID, StartMicroseconds: row.SchoolStartTime.Microseconds, EndMicroseconds: row.SchoolEndTime.Microseconds,
		NumberOfPeriods: row.NumberOfPeriods, PeriodDurationMinutes: row.PeriodDurationMinutes, IntervalDurationMinutes: row.IntervalDurationMinutes,
	}, nil
}
func mapGradeSectionRecord(row db.GradeSection) gradeSectionRecord {
	return gradeSectionRecord{
		ID: row.ID, AcademicYearID: row.AcademicYearID, Name: row.Name,
		IntervalStartMicroseconds: row.IntervalStartTime.Microseconds, IntervalEndMicroseconds: row.IntervalEndTime.Microseconds,
		SectionHeadTeacherID: classroomUUID(row.SectionHeadTeacherID), SortOrder: row.SortOrder,
	}
}
func mapTimetablePeriod(row db.TimetablePeriod) TimetablePeriod {
	return TimetablePeriod{
		ID: row.ID, SortOrder: row.SortOrder, PeriodNumber: classroomInt(row.PeriodNumber),
		StartTime: formatClock(row.StartTime), EndTime: formatClock(row.EndTime), SlotType: row.SlotType,
	}
}

type timetableEntryRepository struct{ queries *db.Queries }

func newTimetableEntryRepository(pool *pgxpool.Pool) *timetableEntryRepository {
	return &timetableEntryRepository{queries: db.New(pool)}
}
func (r *timetableEntryRepository) status(ctx context.Context, id uuid.UUID) (string, error) {
	row, err := r.queries.GetTimetableByID(ctx, id)
	return row.Status, err
}
func (r *timetableEntryRepository) listEntries(ctx context.Context, timetableID uuid.UUID) ([]TimetableEntry, error) {
	rows, err := r.queries.ListTimetableEntriesByTimetable(ctx, timetableID)
	if err != nil {
		return nil, err
	}
	result := make([]TimetableEntry, len(rows))
	for i, row := range rows {
		result[i] = TimetableEntry{
			ID: row.ID, TimetableID: row.TimetableID, DayOfWeek: row.DayOfWeek, PeriodNumber: row.PeriodNumber,
			SubjectID: entryUUID(row.SubjectID), TeacherID: entryUUID(row.TeacherID), ClassroomID: entryUUID(row.ClassroomID),
			CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time, SubjectName: entryText(row.SubjectName),
			TeacherName: entryText(row.TeacherName), ClassroomName: entryText(row.ClassroomName),
		}
	}
	return result, nil
}
func (r *timetableEntryRepository) upsertEntry(ctx context.Context, timetableID uuid.UUID, command timetableEntryCommand) error {
	_, err := r.queries.UpsertTimetableEntry(ctx, db.UpsertTimetableEntryParams{
		TimetableID: timetableID, DayOfWeek: command.DayOfWeek, PeriodNumber: command.PeriodNumber,
		SubjectID: optionalClassroomUUID(command.SubjectID), TeacherID: optionalClassroomUUID(command.TeacherID), ClassroomID: optionalClassroomUUID(command.ClassroomID),
	})
	return err
}
func (r *timetableEntryRepository) deleteEntry(ctx context.Context, timetableID uuid.UUID, day, period int16) error {
	return r.queries.DeleteTimetableEntry(ctx, db.DeleteTimetableEntryParams{TimetableID: timetableID, DayOfWeek: day, PeriodNumber: period})
}

func (r *timetableEntryRepository) validationContext(ctx context.Context, id uuid.UUID) (validationContext, error) {
	timetable, err := r.queries.GetTimetableByID(ctx, id)
	if err != nil {
		return validationContext{}, err
	}
	class, err := r.queries.GetClassByID(ctx, timetable.ClassID)
	if err != nil {
		return validationContext{}, err
	}
	return validationContext{AcademicYearID: timetable.AcademicYearID, ClassID: timetable.ClassID, GradeID: class.GradeID, FormTeacherID: entryUUID(class.FormTeacherID)}, nil
}
func (r *timetableEntryRepository) validationEntries(ctx context.Context, id uuid.UUID) ([]TimetableEntry, error) {
	return r.listEntries(ctx, id)
}
func (r *timetableEntryRepository) crossBookings(ctx context.Context, yearID, timetableID uuid.UUID) ([]crossBooking, error) {
	rows, err := r.queries.ListEntriesForYearExcludingTimetable(ctx, db.ListEntriesForYearExcludingTimetableParams{AcademicYearID: yearID, ID: timetableID})
	if err != nil {
		return nil, err
	}
	result := make([]crossBooking, len(rows))
	for i, row := range rows {
		result[i] = crossBooking{DayOfWeek: row.DayOfWeek, PeriodNumber: row.PeriodNumber, TeacherID: entryUUID(row.TeacherID), ClassroomID: entryUUID(row.ClassroomID), ClassName: row.ClassName}
	}
	return result, nil
}
func (r *timetableEntryRepository) teacherUnavailable(ctx context.Context, teacherID, yearID uuid.UUID, day, period int16) (bool, error) {
	return r.queries.IsTeacherUnavailable(ctx, db.IsTeacherUnavailableParams{TeacherID: teacherID, AcademicYearID: yearID, DayOfWeek: day, PeriodNumber: period})
}
func (r *timetableEntryRepository) classSubjectTeacher(ctx context.Context, classID, subjectID uuid.UUID) (uuid.UUID, error) {
	return r.queries.GetClassSubjectTeacher(ctx, db.GetClassSubjectTeacherParams{ClassID: classID, SubjectID: subjectID})
}
func (r *timetableEntryRepository) teacherAssignedSubject(ctx context.Context, teacherID, subjectID uuid.UUID) (bool, error) {
	return r.queries.IsTeacherAssignedToSubject(ctx, db.IsTeacherAssignedToSubjectParams{TeacherID: teacherID, SubjectID: subjectID})
}
func (r *timetableEntryRepository) requirements(ctx context.Context, yearID, gradeID uuid.UUID) ([]requirementValue, error) {
	rows, err := r.queries.ListSubjectPeriodRequirementsByGrade(ctx, db.ListSubjectPeriodRequirementsByGradeParams{AcademicYearID: yearID, GradeID: gradeID})
	if err != nil {
		return nil, err
	}
	result := make([]requirementValue, len(rows))
	for i, row := range rows {
		result[i] = requirementValue{SubjectID: row.SubjectID, SubjectName: row.SubjectName, PeriodsPerWeek: row.PeriodsPerWeek}
	}
	return result, nil
}
func (r *timetableEntryRepository) entrySubjectCounts(ctx context.Context, timetableID uuid.UUID) (map[uuid.UUID]int32, error) {
	rows, err := r.queries.CountEntriesBySubjectForTimetable(ctx, timetableID)
	if err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID]int32, len(rows))
	for _, row := range rows {
		if row.SubjectID.Valid {
			result[uuid.UUID(row.SubjectID.Bytes)] = row.EntryCount
		}
	}
	return result, nil
}
func (r *timetableEntryRepository) authorizedReviewers(ctx context.Context, yearID, gradeID uuid.UUID) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, 2)
	tic, err := r.queries.GetGradeTICForGrade(ctx, db.GetGradeTICForGradeParams{AcademicYearID: yearID, GradeID: gradeID})
	if err == nil {
		result = append(result, tic)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	section, err := r.queries.GetGradeSectionForGrade(ctx, db.GetGradeSectionForGradeParams{AcademicYearID: yearID, GradeID: gradeID})
	if err == nil && section.SectionHeadTeacherID.Valid {
		result = append(result, uuid.UUID(section.SectionHeadTeacherID.Bytes))
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return result, nil
}

type settingsRepository struct{ queries *db.Queries }

func newSettingsRepository(pool *pgxpool.Pool) *settingsRepository {
	return &settingsRepository{queries: db.New(pool)}
}
func (r *settingsRepository) upsert(ctx context.Context, values settingsValues) (Settings, error) {
	row, err := r.queries.UpsertTimetableSettings(ctx, db.UpsertTimetableSettingsParams{AcademicYearID: values.AcademicYearID, SchoolStartTime: pgtype.Time{Microseconds: values.StartMicroseconds, Valid: true}, SchoolEndTime: pgtype.Time{Microseconds: values.EndMicroseconds, Valid: true}, NumberOfPeriods: values.NumberOfPeriods, PeriodDurationMinutes: values.PeriodDurationMinutes, IntervalDurationMinutes: values.IntervalDurationMinutes})
	return mapSettings(row), err
}
func (r *settingsRepository) getByYear(ctx context.Context, yearID uuid.UUID) (Settings, error) {
	row, err := r.queries.GetTimetableSettingsByYear(ctx, yearID)
	return mapSettings(row), err
}
func formatClock(value pgtype.Time) string {
	if !value.Valid {
		return ""
	}
	seconds := value.Microseconds / 1_000_000
	return fmt.Sprintf("%02d:%02d", seconds/3600, (seconds%3600)/60)
}
func mapSettings(row db.TimetableSetting) Settings {
	return Settings{AcademicYearID: row.AcademicYearID, SchoolStartTime: formatClock(row.SchoolStartTime), SchoolEndTime: formatClock(row.SchoolEndTime), NumberOfPeriods: row.NumberOfPeriods, PeriodDurationMinutes: row.PeriodDurationMinutes, IntervalDurationMinutes: row.IntervalDurationMinutes}
}
