CREATE TABLE organisation_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,

    sub_plan TEXT NOT NULL,
    status TEXT NOT NULL,

    started_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE organisation_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organisation_id UUID NOT NULL
        REFERENCES organisations(id) ON DELETE CASCADE,

    subscription_id UUID
        REFERENCES organisation_subscriptions(id) ON DELETE SET NULL,

    amount NUMERIC(12, 2) NOT NULL,
    currency TEXT NOT NULL,

    status TEXT NOT NULL,

    paid_at TIMESTAMPTZ,
    comment TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organisation_payments_organisation_id
    ON organisation_payments(organisation_id);

CREATE INDEX idx_organisation_payments_subscription_id
    ON organisation_payments(subscription_id);