import { ReactElement, useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';

interface AuthGuardProps {
  children: ReactElement;
  requiredPermission?: string;
}

export function AuthGuard({ children, requiredPermission }: AuthGuardProps): ReactElement {
  const location = useLocation();
  const { isAuthenticated, user, checkAuth } = useAuthStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (requiredPermission && user?.permissions && !user.permissions.includes(requiredPermission)) {
    return <Navigate to="/403" replace />;
  }

  return children;
}