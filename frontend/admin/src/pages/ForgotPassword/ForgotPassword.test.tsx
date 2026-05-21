import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { ForgotPassword } from './index';
import { BrowserRouter } from 'react-router-dom';
import { vi, describe, it, expect } from 'vitest';

const renderForgotPassword = () => {
  render(
    <BrowserRouter>
      <ForgotPassword />
    </BrowserRouter>
  );
};

describe('ForgotPassword', () => {
  it('renders forgot password form with email field', () => {
    renderForgotPassword();
    expect(screen.getByText('忘记密码')).toBeInTheDocument();
    expect(screen.getByText('输入您的注册邮箱，我们会发送重置链接')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('注册邮箱')).toBeInTheDocument();
  });

  it('shows back to login link', () => {
    renderForgotPassword();
    expect(screen.getByText('返回登录')).toBeInTheDocument();
  });

  it('shows send reset button', () => {
    renderForgotPassword();
    expect(screen.getByText('发送重置链接')).toBeInTheDocument();
  });

  it('shows validation error for empty email', async () => {
    renderForgotPassword();
    const button = screen.getByText('发送重置链接');
    fireEvent.click(button);
    await waitFor(() => {
      expect(screen.getByText('请输入邮箱')).toBeInTheDocument();
    });
  });

  it('shows validation error for invalid email format', async () => {
    renderForgotPassword();
    const emailInput = screen.getByPlaceholderText('注册邮箱');
    fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
    const button = screen.getByText('发送重置链接');
    fireEvent.click(button);
    await waitFor(() => {
      expect(screen.getByText('请输入有效的邮箱地址')).toBeInTheDocument();
    });
  });

  it('shows success message after sending reset link', async () => {
    renderForgotPassword();
    const emailInput = screen.getByPlaceholderText('注册邮箱');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    const button = screen.getByText('发送重置链接');

    await act(async () => {
      fireEvent.click(button);
      // Wait for the simulated API call (1500ms) plus some buffer
      await new Promise((resolve) => setTimeout(resolve, 2000));
    });

    expect(screen.getByText('邮件已发送')).toBeInTheDocument();
  });

  it('shows return to login button after success', async () => {
    renderForgotPassword();
    const emailInput = screen.getByPlaceholderText('注册邮箱');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    const button = screen.getByText('发送重置链接');

    await act(async () => {
      fireEvent.click(button);
      await new Promise((resolve) => setTimeout(resolve, 2000));
    });

    expect(screen.getByText('返回登录')).toBeInTheDocument();
  });
});