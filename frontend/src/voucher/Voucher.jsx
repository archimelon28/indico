import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

function GetListVoucher() {
    const [vouchers, setVouchers] = useState([])
    const navigate = useNavigate()
    const [search, setSearch] = useState("")
    // localStorage.removeItem('token')

    useEffect(() => {
        const token = localStorage.getItem('token')
        if (!token) {
            navigate('/login')
        }
        fetch('http://localhost:8080/voucher', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        })
            .then((res) => res.json())
            .then((response) => {
                setVouchers(response.data)
            })
            .catch((err) => {
                console.error('Error fetching data:', err)
            })
    }, [])

    const handleDelete = (id) => {
        const confirmDelete = window.confirm("Yakin ingin menghapus voucher ini?")
        if (!confirmDelete) return

        fetch(`http://localhost:8080/voucher/${id}`, {
            method: 'DELETE',
        })
            .then((res) => {
                if (res.ok) {
                    setVouchers(prev => prev.filter(voucher => voucher.id !== id))
                } else {
                    console.error('Gagal menghapus voucher')
                }
            })
            .catch((err) => {
                console.error('Error saat delete:', err)
            })
    }

    const filteredVouchers = vouchers.filter((v) =>
        v.voucher_code.toLowerCase().includes(search.toLowerCase())
    )

    const handleLogout = () => {
        localStorage.removeItem('token')
        navigate('/login')
    }

    const handleExport = () => {
        fetch("http://localhost:8080/voucher/export", {
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        })
            .then((res) => res.blob())
            .then((blob) => {
                const url = window.URL.createObjectURL(blob)
                const a = document.createElement("a")
                a.href = url
                a.download = "vouchers.csv"
                a.click()
            })
    }

    return (
        <div>
            <h1>Voucher List</h1>
            <input
                type="text"
                placeholder="Cari Voucher Code..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
            />
            <br />
            <button type="button" onClick={() => navigate('/voucher')}>
                Insert New Voucher
            </button>
            <br />
            <button type="button" onClick={() => navigate('/upload-csv')}>
                Upload CSV
            </button>
            <table border="1" cellPadding="10">
                <thead>
                    <tr>
                        <th>No</th>
                        <th>Voucher Code</th>
                        <th>Discount (%)</th>
                        <th>Expiry Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredVouchers.map((voucher, index) => (
                        <tr key={voucher.id}>
                            <td>{index + 1}</td>
                            <td>{voucher.voucher_code}</td>
                            <td>{voucher.discount_percent}%</td>
                            <td>
                                {new Intl.DateTimeFormat('id-ID').format(new Date(voucher.expiry_date))}
                            </td>
                            <td>
                                <button onClick={() => navigate(`/voucher/${voucher.id}`)}>Edit</button>
                                <button onClick={() => handleDelete(voucher.id)}>Delete</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <button onClick={handleExport}>Export CSV</button>
            <br />
            <button type="button" onClick={handleLogout}>
                Logout
            </button>
        </div>
    )
}

export default GetListVoucher
