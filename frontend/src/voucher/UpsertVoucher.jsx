import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

function InsertVoucher() {
    const navigate = useNavigate()
    const { id } = useParams()
    const [form, setForm] = useState({
        voucher_code: '',
        discount_percent: '',
        expiry_date: '',
    })

    const isEdit = !!id

    useEffect(() => {
        const token = localStorage.getItem('token')
        if (!token) {
            navigate('/login')
        }
        if (isEdit) {
            fetch(`http://localhost:8080/voucher/${id}`)
                .then((res) => res.json())
                .then((data) => {
                    setForm({
                        voucher_code: data.voucher_code,
                        discount_percent: data.discount_percent,
                        expiry_date: data.expiry_date.slice(0, 10),
                    })
                })
        }
    }, [id])

    const handleChange = (e) => {
        const { name, value } = e.target
        setForm({
            ...form,
            [name]: name === 'discount_percent' ? parseInt(value) : value

        })
    }

    const handleSubmit = (e) => {
        e.preventDefault()
        if (!form.voucher_code || !form.discount_percent || !form.expiry_date) {
            alert('Semua field harus diisi!')
            return
        }

        if (form.discount_percent < 1 || form.discount_percent > 100) {
            alert('Discount harus antara 1 - 100')
            return
        }
        const method = isEdit ? 'PUT' : 'POST'
        const url = isEdit
            ? `http://localhost:8080/voucher/${id}`
            : 'http://localhost:8080/voucher'

        fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(form),
        })
            .then((res) => {
                if (res.ok) {
                    navigate('/')
                } else {
                    console.error('Gagal simpan data')
                }
            })
            .catch((err) => console.error('Error:', err))
    }

    return (
        <div>
            <h1>{isEdit ? 'Edit' : 'Insert'} Voucher</h1>
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    name="voucher_code"
                    value={form.voucher_code}
                    onChange={handleChange}
                    placeholder="Voucher Code"
                />
                <input
                    type="number"
                    name="discount_percent"
                    value={form.discount_percent}
                    onChange={handleChange}
                    placeholder="Discount (%)"
                />
                <input
                    type="date"
                    name="expiry_date"
                    value={form.expiry_date}
                    onChange={handleChange}
                />
                <button type="submit">
                    {isEdit ? 'Update' : 'Insert'}
                </button>
                <button type="button" onClick={() => navigate('/')}>
                    Cancel
                </button>
            </form>
        </div>
    )
}

export default InsertVoucher