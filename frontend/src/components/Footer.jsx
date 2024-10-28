import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/footer.scss';

function Footer() {
    return (
        <footer className="custom-footer py-4">
            <div className="container text-center text-md-start">
                <div className="row footer-cust">
                    {/* Footer Logo and Description */}
                    <div className="col-md-4 mb-5">
                        <h5 className="text-uppercase">DocuDefense</h5>
                        <p>Your trusted partner in document security and management.</p>
                    </div>

                    {/* Footer Contact Information */}
                    <div className="col-md-4 mb-4">
                        <h5 className="text-uppercase">Contact</h5>
                        <p>Email: <a href="mailto:support@docudefense.com" className="footer-link">support@docudefense.com</a></p>
                        <p>Phone: <a href="tel:+1234567890" className="footer-link">+1 234 567 890</a></p>
                    </div>
                </div>

                {/* Copyright and Additional Links */}
                <div className="text-center pt-3">
                    <p>© {new Date().getFullYear()} DocuDefense. All rights reserved.</p>
                    <p>
                        Need help or have questions? <Link to="/contact" className="footer-link">Contact Support</Link>
                    </p>
                </div>
            </div>
        </footer>
    );
}

export default Footer;
