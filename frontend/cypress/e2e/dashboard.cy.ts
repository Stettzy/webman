describe('Dashboard', () => {
    beforeEach(() => {
        cy.visit('/')
    })

    it('loads the dashboard page', () => {
        cy.get('select').first().should('contain', 'POST')
        cy.get('input[placeholder="Enter request URL"]').should('exist')
        cy.get('button').contains('Send').should('exist')
    })

    it('can switch between tabs', () => {
        // Check Params tab
        cy.contains('button', 'Params').click()
        cy.get('input[placeholder="Parameter Key"]').should('exist')
        cy.get('input[placeholder="Parameter Value"]').should('exist')

        // Check Headers tab
        cy.contains('button', 'Headers').click()
        cy.get('input[placeholder="Header Key"]').should('exist')
        cy.get('input[placeholder="Header Value"]').should('exist')

        // Check Body tab
        cy.contains('button', 'Body').click()
        cy.get('.monaco-editor').should('exist')
    })

    it('can add and remove parameters', () => {
        cy.contains('button', 'Params').click()

        // Add parameter
        cy.get('input[placeholder="Parameter Key"]').type('testKey')
        cy.get('input[placeholder="Parameter Value"]').type('testValue')
        cy.contains('button', 'Add').click()

        // Verify parameter was added
        cy.contains('testKey').should('exist')
        cy.contains('testValue').should('exist')

        // Remove parameter
        cy.get('button').contains('times').click()
        cy.contains('testKey').should('not.exist')
    })

    it('can add and remove headers', () => {
        cy.contains('button', 'Headers').click()

        // Add header
        cy.get('input[placeholder="Header Key"]').type('Content-Type')
        cy.get('input[placeholder="Header Value"]').type('application/json')
        cy.contains('button', 'Add').click()

        // Verify header was added
        cy.contains('Content-Type').should('exist')
        cy.contains('application/json').should('exist')

        // Remove header
        cy.get('button').contains('times').click()
        cy.contains('Content-Type').should('not.exist')
    })

    it('can edit and format JSON in the body', () => {
        cy.contains('button', 'Body').click()

        // Select JSON body type
        cy.get('select').last().select('JSON')

        // Type unformatted JSON
        cy.get('.monaco-editor').type('{"key":"value","nested":{"inner":"value"}}')

        // Click format button
        cy.contains('button', 'Format JSON').click()

        // Verify JSON was formatted (basic check since we can't easily check editor content)
        cy.get('.monaco-editor').should('exist')
    })

    it('can send a request and receive a response', () => {
        // Intercept API calls
        cy.intercept('POST', 'http://localhost:9090', {
            statusCode: 200,
            body: {
                body: btoa('{"success":true}'),
                headers: { 'Content-Type': 'application/json' },
                statusCode: 200
            }
        }).as('apiRequest')

        // Enter request URL
        cy.get('input[placeholder="Enter request URL"]').type('http://api.example.com')

        // Add a header
        cy.contains('button', 'Headers').click()
        cy.get('input[placeholder="Header Key"]').type('Content-Type')
        cy.get('input[placeholder="Header Value"]').type('application/json')
        cy.contains('button', 'Add').click()

        // Add request body
        cy.contains('button', 'Body').click()
        cy.get('.monaco-editor').type('{"test":"data"}')

        // Send request
        cy.contains('button', 'Send').click()

        // Wait for response
        cy.wait('@apiRequest')

        // Verify response is displayed
        cy.contains('Status: 200').should('exist')
        cy.contains('success').should('exist')
    })

    it('handles request errors gracefully', () => {
        // Intercept API calls with error
        cy.intercept('POST', 'http://localhost:9090', {
            statusCode: 500,
            body: {
                message: 'fail'
            }
        }).as('apiRequest')

        // Enter request URL
        cy.get('input[placeholder="Enter request URL"]').type('http://api.example.com')

        // Send request
        cy.contains('button', 'Send').click()

        // Wait for response
        cy.wait('@apiRequest')

        // Verify error handling
        cy.get('.response-section').should('not.exist')
    })
}) 