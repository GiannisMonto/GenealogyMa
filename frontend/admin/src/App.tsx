import { ReactElement, Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import { Spin } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AppLayout } from '@/components/Layout';
import { AuthGuard } from '@/components/AuthGuard';
import { Dashboard } from '@/pages/Dashboard';
import { Login } from '@/pages/Login';
import { Members } from '@/pages/Members';
import { MembersDetail } from '@/pages/Members/Detail';
import { Users } from '@/pages/Users';
import { Roles } from '@/pages/Roles';
import { Audit } from '@/pages/Audit';
import { NotFound } from '@/pages/Error/404';

const loading = (
  <Spin
    size="large"
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
    }}
  />
);

const App: React.FC = () => (
  <ConfigProvider locale={zhCN}>
    <BrowserRouter>
      <Suspense fallback={loading}>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <AuthGuard>
                <AppLayout />
              </AuthGuard>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="members" element={<Members />} />
            <Route path="members/:id" element={<MembersDetail />} />
            <Route path="users" element={<Users />} />
            <Route path="roles" element={<Roles />} />
            <Route path="audit" element={<Audit />} />
          </Route>
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  </ConfigProvider>
);

export default App;