package logwriter_test

import (
	"database/sql"

	testdb "github.com/erikstmartin/go-testdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/galera-healthcheck/cluster-health-logger/logwriter"
	"os"
	"io/ioutil"
)

var (
	logFile *os.File
)

var _ = Describe("Cluster Health Logger", func() {

	BeforeEach(func() {
		logFile, _ = ioutil.TempFile(os.TempDir(), "logFile")

	})

	AfterEach(func() {
		os.Remove(logFile.Name())
	})

	Context("when the log file does not exist", func() {
		BeforeEach(func() {
			os.Remove(logFile.Name())
		})

		It("writes headers to the file", func() {
			logWriter := logWriterTestHelper(logFile.Name())
			ts := "happy-time"
			logWriter.Write(ts)
			contents, err := ioutil.ReadFile(logFile.Name())
			Expect(err).ToNot(HaveOccurred())
			contentsStr := string(contents)
			Expect(contentsStr).To(Equal("timestamp,a,b,c,d,e,f,g,h,i\nhappy-time,1,2,3,4,5,6,7,8,9\n"))
		})
	})

	Context("when the log file exists", func() {

		It("writes only the rows to the file", func() {
			logWriter := logWriterTestHelper(logFile.Name())
			ts := "happy-time"
			logWriter.Write(ts)
			contents, err := ioutil.ReadFile(logFile.Name())
			Expect(err).ToNot(HaveOccurred())
			contentsStr := string(contents)
			Expect(contentsStr).To(Equal("happy-time,1,2,3,4,5,6,7,8,9\n"))
		})

		It("writes a new line", func() {
			logWriter := logWriterTestHelper(logFile.Name())
			ts1 := "happy-time"
			logWriter.Write(ts1)
			ts2 := "sad-time"
			logWriter.Write(ts2)
			contents, err := ioutil.ReadFile(logFile.Name())
			Expect(err).ToNot(HaveOccurred())
			contentsStr := string(contents)
			Expect(contentsStr).To(Equal("happy-time,1,2,3,4,5,6,7,8,9\nsad-time,1,2,3,4,5,6,7,8,9\n"))
		})

	})
})

func logWriterTestHelper(filePath string) logwriter.LogWriter {
	db, _ := sql.Open("testdb", "")

	sql := "SHOW STATUS WHERE Variable_name IN ('wsrep_ready','wsrep_cluster_conf_id','wsrep_cluster_status','wsrep_connected','wsrep_local_state_comment','wsrep_local_recv_queue_avg','wsrep_flow_control_paused','wsrep_cert_deps_distance','wsrep_local_send_queue_avg')"
	columns := []string{"Variable_name", "Value"}
	result := "a,1\nb,2\nc,3\nd,4\ne,5\nf,6\ng,7\nh,8\ni,9"
	testdb.StubQuery(sql, testdb.RowsFromCSVString(columns, result))

	return logwriter.New(db, filePath)
}