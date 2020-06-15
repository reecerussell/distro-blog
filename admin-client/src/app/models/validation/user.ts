const rules = {
    firstname: {
        required: {
            value: true,
            error: "The firstname field is required.",
        },
        maxLength: {
            value: 45,
            error: "Firstname cannot be greater than 45 characters long.",
        },
    },
    lastname: {
        required: {
            value: true,
            error: "The lastname field is required.",
        },
        maxLength: {
            value: 45,
            error: "Lastname cannot be greater than 45 characters long.",
        },
    },
    email: {
        required: {
            value: true,
            error: "The email field is required.",
        },
        maxLength: {
            value: 100,
            error: "Email cannot be greater than 100 characters long.",
        },
    },
};

const Validate = (model: any): any => {
    const keys = Object.keys(model);

    for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        const value = model[key];

        const ruleSet = rules[key];
        if (!ruleSet) {
            continue;
        }

        if (typeof value === "string") {
            if ((!value || value.length < 1) && ruleSet.required.value) {
                return ruleSet.required.error;
            }

            if (value.length > ruleSet.maxLength.value) {
                return ruleSet.maxLength.error;
            }
        }
    }
};

export default Validate;
