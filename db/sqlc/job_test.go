package db

import (
	"context"
	"database/sql"
	"example/employee/server/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomJob(t *testing.T) Job {
	arg := CreateJobParams{
		JobTitle:  util.RandomJobTitle(),
		MinSalary: sql.NullInt64{util.RandomMinSalary(), true},
		MaxSalary: sql.NullInt64{util.RandomMaxSalary(), true},
	}
	job, err := testQueries.CreateJob(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, job)

	require.Equal(t, job.JobTitle, arg.JobTitle)
	require.Equal(t, job.MinSalary, arg.MinSalary)
	require.Equal(t, job.MaxSalary, arg.MaxSalary)
	require.True(t, job.MinSalary.Int64 <= job.MaxSalary.Int64)

	require.NotZero(t, job.JobID)

	return job
}

func TestCreateJob(t *testing.T) {
	createRandomJob(t)
}

func TestGetJob(t *testing.T) {
	job1 := createRandomJob(t)
	job2, err := testQueries.GetJob(context.Background(), job1.JobID)

	require.NoError(t, err)
	require.NotEmpty(t, job2)

	require.Equal(t, job1.JobTitle, job2.JobTitle)
	require.Equal(t, job1.MinSalary, job2.MinSalary)
	require.Equal(t, job1.MaxSalary, job2.MaxSalary)
}

func TestUpdateJob(t *testing.T) {
	job1 := createRandomJob(t)
	updateJobTitle := util.RandomJobTitle()
	arg := UpdateJobParams{
		JobID:     job1.JobID,
		JobTitle:  updateJobTitle,
		MinSalary: job1.MinSalary,
		MaxSalary: job1.MaxSalary,
	}

	job2, err := testQueries.UpdateJob(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, job2)

	require.Equal(t, job1.JobID, job2.JobID)
	require.Equal(t, arg.JobTitle, job2.JobTitle)
	require.Equal(t, job1.MinSalary, job2.MinSalary)
	require.Equal(t, job1.MaxSalary, job2.MaxSalary)
}

func TestDeleteJob(t *testing.T) {
	job1 := createRandomJob(t)
	err := testQueries.DeleteJob(context.Background(), job1.JobID)
	require.NoError(t, err)

	job2, err := testQueries.GetJob(context.Background(), job1.JobID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, job2)
}

func TestListJob(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomJob(t)
	}

	arg := ListJobsParams{
		Limit:  5,
		Offset: 5,
	}

	jobs, err := testQueries.ListJobs(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, jobs, 5)

	for _, job := range jobs {
		require.NotEmpty(t, job)
	}
}
