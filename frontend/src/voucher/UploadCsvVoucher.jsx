import React, { useState } from 'react'

function UploadCSV() {
  const [file, setFile] = useState(null)
  const [preview, setPreview] = useState([])
  const [result, setResult] = useState(null)

  const handleFileChange = (e) => {
    const csvFile = e.target.files[0]
    setFile(csvFile)

    const reader = new FileReader()
    reader.onload = (e) => {
      const lines = e.target.result.split("\n").slice(1)
      const data = lines.map(line => line.split(","))
      setPreview(data)
    }
    reader.readAsText(csvFile)
  }

  const handleUpload = () => {
    const formData = new FormData()
    formData.append('file', file)

    fetch('http://localhost:8080/voucher/upload-csv', {
      method: 'POST',
      headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
      body: formData,
    })
      .then((res) => res.json())
      .then(setResult)
  }

  return (
    <div>
      <h2>Upload CSV Voucher</h2>
      <input type="file" accept=".csv" onChange={handleFileChange} />
      {preview.length > 0 && (
        <div>
          <h4>Preview</h4>
          <table>
            <thead>
              <tr><th>Kode</th><th>Diskon</th><th>Expiry</th></tr>
            </thead>
            <tbody>
              {preview.map((row, i) => (
                <tr key={i}>
                  <td>{row[0]}</td><td>{row[1]}</td><td>{row[2]}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <button onClick={handleUpload}>Upload</button>
        </div>
      )}
      {result && (
        <p>{result.success} berhasil, {result.failed} gagal</p>
      )}
    </div>
  )
}

export default UploadCSV
