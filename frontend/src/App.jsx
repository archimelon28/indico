import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import GetListVoucher from './voucher/Voucher'
import UpsertVoucher from './voucher/UpsertVoucher'
import Login from './login/Login'
import UploadCSV from './voucher/UploadCsvVoucher'
import './App.css'

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<GetListVoucher />} />
                <Route path="/login" element={<Login />} />
                <Route path="/upload-csv" element={<UploadCSV />} />
                <Route path="/voucher" element={<UpsertVoucher />} />
                <Route path="/voucher/:id" element={<UpsertVoucher />} /> 
            </Routes>
        </Router>
    )
}

export default App