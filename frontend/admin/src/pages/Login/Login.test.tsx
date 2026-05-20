import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Login } from './index';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect, beforeEach } from 'vitest';

const mockLogin = vi.fn();
vi.mock('@/store/auth', () => ({
  useAuthStore: () => ({
    login: mockLogin,
    isLoading: false,
  }),
}));

const renderLogin = () => {
  render(
    <BrowserRouter>
      <Login />
    </BrowserRouter>
  );
};

describe('Login', () => {
  beforeEach(() => {
    mockLogin.mockReset();
  });

  it('renders login form with all required fields', () => {
    renderLogin();
    expect(screen.getByPlaceholderText('用户名')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('密码')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('验证码')).toBeInTheDocument();
  });

  it('shows forgot password link', () => {
    renderLogin();
    expect(screen.getByText('忘记密码？')).toBeInTheDocument();
  });

  it('shows remember me checkbox', () => {
    renderLogin();
    expect(screen.getByText('记住我')).toBeInTheDocument();
  });

  it('displays captcha image', () => {
    renderLogin();
    const captchaImg = screen.getByAltText('验证码');
    expect(captchaImg).toBeInTheDocument();
  });
});